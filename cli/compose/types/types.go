package types

// ServiceVolumeConfig are references to a volume used by a service
type ServiceVolumeConfig struct {
	Type        string
	Source      string
	Target      string
	ReadOnly    bool `mapstructure:"read_only"`
	Consistency string
	Bind        *ServiceVolumeBind
	Volume      *ServiceVolumeVolume
}

// ServiceVolumeBind are options for a service volume of type bind
type ServiceVolumeBind struct {
	Propagation string
}

// ServiceVolumeVolume are options for a service volume of type volume
type ServiceVolumeVolume struct {
	NoCopy bool `mapstructure:"nocopy"`
}
