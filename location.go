package ark

type Location struct {
	X float32
	Y float32
	Z float32

	Pitch float32
	Yaw   float32
	Roll  float32
}

func readLocation(vr valueReader) (Location, error) {
	var err error
	var l Location

	l.X, err = vr.readFloat()
	if err != nil {
		return l, err
	}

	l.Y, err = vr.readFloat()
	if err != nil {
		return l, err
	}

	l.Z, err = vr.readFloat()
	if err != nil {
		return l, err
	}

	l.Pitch, err = vr.readFloat()
	if err != nil {
		return l, err
	}

	l.Yaw, err = vr.readFloat()
	if err != nil {
		return l, err
	}

	l.Roll, err = vr.readFloat()
	if err != nil {
		return l, err
	}

	return l, nil
}
