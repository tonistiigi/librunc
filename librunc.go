package librunc

// Takes in empty directory and json configuration
func New(dir, config string) (*Container, error) {

}

// Takes in directory. Tries to open configuration from that directory.
func NewFromDirectory(dir string) error {

}

func (*Container) Start() error { // todo: fds

}

func (*Container) Kill() error {

}

// func (*Container) Stats(interval int)
// func (*Container) NotifyOOM()
// func (*Container) Checkpoint()
// func (*Container) Restore()
