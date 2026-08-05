package Utils

type Scale struct {
	Rate float64

	//The width / height of original image. These are not expected to reassign
	OriginalWidth  float64
	OriginalHeight float64
}

func (s *Scale) Size() (float64, float64) {
	return s.OriginalWidth * s.Rate, s.OriginalHeight * s.Rate
}
