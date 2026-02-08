package companion

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/langowarny/lango/internal/logging"
)

var logger = logging.SubsystemSugar("companion.discovery")

const (
	// ServiceType is the mDNS service type for Lango companion apps.
	ServiceType = "_lango-companion._tcp"
	// Domain is the mDNS domain.
	Domain = "local."
)

// ServiceInfo represents a discovered companion service.
type ServiceInfo struct {
	Name     string
	Host     string
	Port     int
	IPs      []net.IP
	TXT      map[string]string
	LastSeen time.Time
}

// Address returns the WebSocket address for the service.
func (s *ServiceInfo) Address() string {
	if len(s.IPs) > 0 {
		return fmt.Sprintf("ws://%s:%d", s.IPs[0].String(), s.Port)
	}
	return fmt.Sprintf("ws://%s:%d", s.Host, s.Port)
}

// Discovery handles mDNS service discovery for companion apps.
type Discovery struct {
	mu        sync.RWMutex
	services  map[string]*ServiceInfo
	resolver  *zeroconf.Resolver
	callbacks []func(*ServiceInfo)
	running   bool
	stopCh    chan struct{}
}

// NewDiscovery creates a new Discovery instance.
func NewDiscovery() *Discovery {
	return &Discovery{
		services: make(map[string]*ServiceInfo),
		stopCh:   make(chan struct{}),
	}
}

// Start begins service discovery.
func (d *Discovery) Start(ctx context.Context) error {
	d.mu.Lock()
	if d.running {
		d.mu.Unlock()
		return nil
	}
	d.running = true
	d.mu.Unlock()

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return fmt.Errorf("failed to create resolver: %w", err)
	}
	d.resolver = resolver

	entries := make(chan *zeroconf.ServiceEntry)

	go func() {
		for {
			select {
			case entry := <-entries:
				if entry != nil {
					d.handleEntry(entry)
				}
			case <-d.stopCh:
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	logger.Infow("starting companion discovery", "serviceType", ServiceType)

	go func() {
		for {
			select {
			case <-d.stopCh:
				return
			case <-ctx.Done():
				return
			default:
				browseCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				err := d.resolver.Browse(browseCtx, ServiceType, Domain, entries)
				cancel()
				if err != nil {
					logger.Warnw("browse error", "error", err)
				}
				time.Sleep(30 * time.Second)
			}
		}
	}()

	return nil
}

// Stop stops service discovery.
func (d *Discovery) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.running {
		return
	}

	close(d.stopCh)
	d.running = false
	logger.Infow("stopped companion discovery")
}

// GetServices returns all discovered services.
func (d *Discovery) GetServices() []*ServiceInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()

	result := make([]*ServiceInfo, 0, len(d.services))
	for _, s := range d.services {
		result = append(result, s)
	}
	return result
}

// GetService returns a specific service by name.
func (d *Discovery) GetService(name string) *ServiceInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.services[name]
}

// OnServiceFound registers a callback for new service discoveries.
func (d *Discovery) OnServiceFound(callback func(*ServiceInfo)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.callbacks = append(d.callbacks, callback)
}

func (d *Discovery) handleEntry(entry *zeroconf.ServiceEntry) {
	d.mu.Lock()
	defer d.mu.Unlock()

	txt := make(map[string]string)
	for _, t := range entry.Text {
		txt[t] = t
	}

	info := &ServiceInfo{
		Name:     entry.ServiceInstanceName(),
		Host:     entry.HostName,
		Port:     entry.Port,
		IPs:      append(entry.AddrIPv4, entry.AddrIPv6...),
		TXT:      txt,
		LastSeen: time.Now(),
	}

	_, exists := d.services[info.Name]
	d.services[info.Name] = info

	if !exists {
		logger.Infow("discovered companion",
			"name", info.Name,
			"host", info.Host,
			"port", info.Port,
			"ips", info.IPs,
		)

		// Notify callbacks
		for _, cb := range d.callbacks {
			go cb(info)
		}
	}
}

// Config holds companion configuration.
type Config struct {
	// ManualAddress overrides discovery with a specific address.
	ManualAddress string
	// Enabled controls whether companion features are enabled.
	Enabled bool
}

// Manager manages companion connections.
type Manager struct {
	discovery *Discovery
	config    Config
	connected bool
	mu        sync.RWMutex
}

// NewManager creates a new companion Manager.
func NewManager(cfg Config) *Manager {
	return &Manager{
		discovery: NewDiscovery(),
		config:    cfg,
	}
}

// Start starts the companion manager.
func (m *Manager) Start(ctx context.Context) error {
	if !m.config.Enabled {
		logger.Infow("companion features disabled")
		return nil
	}

	// If manual address is configured, use that
	if m.config.ManualAddress != "" {
		logger.Infow("using manual companion address", "address", m.config.ManualAddress)
		// TODO: Connect to manual address
		return nil
	}

	// Otherwise, start discovery
	m.discovery.OnServiceFound(func(info *ServiceInfo) {
		m.handleDiscovery(info)
	})

	return m.discovery.Start(ctx)
}

// Stop stops the companion manager.
func (m *Manager) Stop() {
	m.discovery.Stop()
}

// IsConnected returns true if connected to a companion.
func (m *Manager) IsConnected() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connected
}

func (m *Manager) handleDiscovery(info *ServiceInfo) {
	logger.Infow("attempting to connect to companion", "name", info.Name, "address", info.Address())
	// TODO: Establish WebSocket connection
}
