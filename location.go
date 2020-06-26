package main

type Location struct {
	X float32
	Y float32
	Z float32

	Pitch float32
	Yaw   float32
	Roll  float32
}

func readLocation(a *Archive) (Location, error) {
	var err error
	var l Location

	l.X, err = a.readFloat()
	if err != nil {
		return l, err
	}

	l.Y, err = a.readFloat()
	if err != nil {
		return l, err
	}

	l.Z, err = a.readFloat()
	if err != nil {
		return l, err
	}

	l.Pitch, err = a.readFloat()
	if err != nil {
		return l, err
	}

	l.Yaw, err = a.readFloat()
	if err != nil {
		return l, err
	}

	l.Roll, err = a.readFloat()
	if err != nil {
		return l, err
	}

	return l, nil
}
