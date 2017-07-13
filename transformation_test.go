package cloudinary

import "testing"

func TestTransformation_String(t *testing.T) {
	testCases := []struct {
		transformations *Transformation
		str             string
	}{
		{
			&Transformation{
				Quality: QualityAutoEco,
			},
			"q_auto:eco",
		},
		{
			nil,
			"",
		},
		{
			&Transformation{
				Quality: QualityAuto,
				Width:   123,
				Height:  123,
			},
			"q_auto/w_123,h_123",
		},
		{
			&Transformation{
				Quality: QualityAuto,
				Width:   0.5,
			},
			"q_auto/w_0.5",
		},
		{
			&Transformation{
				Quality: QualityAuto,
				Height:   0.5,
			},
			"q_auto/h_0.5",
		},
		{
			&Transformation{
				Quality: QualityAuto,
				Width:   123,
			},
			"q_auto/w_123",
		},
		{
			&Transformation{
				Quality: QualityAuto,
				Height:  123,
			},
			"q_auto/h_123",
		},
		{
			&Transformation{
				Quality:  QualityAuto,
				Width:    123,
				Height:   123,
				CropMode: CropModeFit,
			},
			"q_auto/w_123,h_123,c_fit",
		},
	}

	for _, tc := range testCases {
		str := tc.transformations.String()
		if str != tc.str {
			t.Errorf("Result string must be equal [%s]. Got: [%s]", tc.str, str)
		}
	}
}
