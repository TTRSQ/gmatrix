package gmatrix

import "errors"

type Matrix struct {
	rowNum int
	colNum int
	datas  []float64
}

func NewMatrix(r, c int, datas []float64) (*Matrix, error) {
	if r <= 0 {
		return nil, errors.New("invalid row length")
	}
	if c <= 0 {
		return nil, errors.New("invalid col length")
	}

	if len(datas) != r*c {
		return nil, errors.New("invalid data length")
	}

	return &Matrix{
		rowNum: r,
		colNum: c,
		datas:  datas,
	}, nil
}

func (ma *Matrix) Add(mb *Matrix) error {
	isSame := ma.sameShape(mb)
	if !isSame {
		return errors.New("matrix different shape")
	}
	for i := range ma.datas {
		ma.datas[i] += mb.datas[i]
	}
	return nil
}

func (ma *Matrix) Sub(mb *Matrix) error {
	isSame := ma.sameShape(mb)
	if !isSame {
		return errors.New("matrix different shape")
	}
	for i := range ma.datas {
		ma.datas[i] -= mb.datas[i]
	}
	return nil
}

func (ma *Matrix) Mul(mb *Matrix) error {
	if ma.colNum != mb.rowNum {
		return errors.New("invalid shape")
	}
	newR := ma.rowNum
	newC := mb.colNum
	newDatas := []float64{}

	common := ma.colNum

	for r := 0; r < newR; r++ {
		for c := 0; c < newC; c++ {
			sum := 0.0
			for com := 0; com < common; com++ {
				sum += ma.datas[r*common+com] * mb.datas[com*newC+c]
			}
			newDatas = append(newDatas, sum)
		}
	}

	ma.rowNum = newR
	ma.colNum = newC
	ma.datas = newDatas

	return nil
}

func (m *Matrix) Func(f func(float64) (float64, error)) error {
	for i := range m.datas {
		val, err := f(m.datas[i])
		if err != nil {
			return err
		}
		m.datas[i] = val
	}
	return nil
}

func (ma *Matrix) sameShape(mb *Matrix) bool {
	return ma.colNum == mb.colNum && ma.rowNum == mb.rowNum
}
