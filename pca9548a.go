package pca9548a

import (
	"fmt"
	"log"
	"sync"

	"golang.org/x/exp/io/i2c"
)

const (
	defaultAddress = 0x70
)

type PCA9548A struct {
	address uint8
	dev     *i2c.Device
	// Mutex protects port from changing during a read / write
	sync.Mutex
	port uint8
}

// NewMux create a new multiplexer on a i2c bus
func NewMux(i2cbus string, opts ...func(*PCA9548A) error) (*PCA9548A, error) {
	pca := new(PCA9548A)
	pca.address = defaultAddress
	pca.port = 127

	for _, option := range opts {
		option(pca)
	}

	var err error
	pca.dev, err = i2c.Open(&i2c.Devfs{Dev: i2cbus}, int(pca.address))
	if err != nil {
		log.Panic(err)
	}

	return pca, nil
}

// Close closes the Mux
func (p *PCA9548A) Close() error {
	return p.dev.Close()
}

// Address sets the i2c address in not using the default address of 0x40
func Address(address uint8) func(*PCA9548A) error {
	return func(p *PCA9548A) error {
		return p.setAddress(address)
	}
}
func (p *PCA9548A) setAddress(address uint8) error {
	p.address = address
	return nil
}

// SetPort switches the multiplexer to desired port
func (p *PCA9548A) SetPort(port uint8) error {
	p.Lock()
	if p.port == port {
		p.Unlock()
		return nil
	}

	if port < 0 || port > 7 {
		p.Unlock()
		return fmt.Errorf("error setting port to %d : port must be be 0-7", port)
	}
	if err := p.dev.Write([]byte{byte(port)}); err != nil {
		p.Unlock()
		return err
	}
	p.port = port
	p.Unlock()
	return nil
}
