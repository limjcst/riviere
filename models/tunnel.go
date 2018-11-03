package models

// Tunnel is the content of TunnelBody
// swagger:model
type Tunnel struct {
	// A port of the gate
	//
	// required: true
	Port int `json:"port"`
	// The address of the target host
	//
	// required: true
	ForwardAddress string `json:"forward_address"`
	// The port of the target host
	//
	// required: true
	ForwardPort int `json:"forward_port"`
}

// PrepareTunnel setups table for tunnel
func PrepareTunnel(db *DB) error {
	q := `
CREATE TABLE IF NOT EXISTS tunnel (port INTEGER PRIMARY KEY,
                                   forward_port INTEGER,
                                   forward_address VARCHAR)
    `
	_, err := db.Exec(q)
	return err
}

// NewTunnel inserts a new record of tunnel
func (db *DB) NewTunnel(tunnel *Tunnel) error {
	q := `
INSERT INTO tunnel (port, forward_port, forward_address) VALUES(?, ?, ?);
    `
	_, err := db.Exec(q, tunnel.Port, tunnel.ForwardPort, tunnel.ForwardAddress)
	return err
}

// DeleteTunnel deletes a record of tunnel based on port
func (db *DB) DeleteTunnel(tunnel *Tunnel) (int64, error) {
	q := `
DELETE FROM tunnel WHERE port=?;
    `
	result, err := db.Exec(q, tunnel.Port)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// ListTunnel lists all the tunnels
func (db *DB) ListTunnel() ([]*Tunnel, error) {
	q := `SELECT port, forward_port, forward_address FROM tunnel;`
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tunnels := make([]*Tunnel, 0)
	for rows.Next() {
		tunnel := new(Tunnel)
		err := rows.Scan(&tunnel.Port, &tunnel.ForwardPort, &tunnel.ForwardAddress)
		if err != nil {
			return nil, err
		}
		tunnels = append(tunnels, tunnel)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tunnels, nil
}
