package detectors

import "github.com/VineBalloon/nozobot/client"

// Detector is the interface for message detectors.
// Message detectors bypass the prefix requirement.
// Note: MsgDetect should expect bad input and error only on internal errors.
// e.g. User sends an invalid format -> return nil.
// Detector fails http request -> return err
type Detector interface {
	Name() string                        /* Returns name of detector */
	Desc() string                        /* Returns name of detector */
	MsgDetect(*client.ClientState) error /* Handles Message Detection */
	Apiget() error                       /* Gets any required API keys */
}
