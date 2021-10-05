package gmatrix

import (
	"errors"
	"math/rand"
)

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

func (m *Matrix) R() int {
	return m.rowNum
}

func (m *Matrix) C() int {
	return m.colNum
}

func (ma *Matrix) Add(mb *Matrix) (*Matrix, error) {
	isSame := ma.sameShape(mb)
	if !isSame {
		return nil, errors.New("matrix different shape")
	}
	datas := []float64{}
	for i := range ma.datas {
		datas = append(datas, ma.datas[i]+mb.datas[i])
	}
	return NewMatrix(ma.rowNum, ma.colNum, datas)
}

func (ma *Matrix) Sub(mb *Matrix) (*Matrix, error) {
	isSame := ma.sameShape(mb)
	if !isSame {
		return nil, errors.New("matrix different shape")
	}
	datas := []float64{}
	for i := range ma.datas {
		datas = append(datas, ma.datas[i]-mb.datas[i])
	}
	return NewMatrix(ma.rowNum, ma.colNum, datas)
}

func (ma *Matrix) Mul(mb *Matrix) (*Matrix, error) {
	if ma.colNum != mb.rowNum {
		return nil, errors.New("invalid shape")
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

	return NewMatrix(newR, newC, newDatas)
}

func (ma *Matrix) MulParallel(mb *Matrix) (*Matrix, error) {
	if ma.colNum != mb.rowNum {
		return nil, errors.New("invalid shape")
	}
	newR := ma.rowNum
	newC := mb.colNum
	newDatas := make([]float64, newR*newC)

	common := ma.colNum
	type ret struct {
		idx int
		val float64
	}
	channel := make(chan ret, newR*newC)

	for r := 0; r < newR; r++ {
		for c := 0; c < newC; c++ {
			go func(row, col int, ch chan ret) {
				sum := 0.0
				for com := 0; com < common; com++ {
					sum += ma.datas[row*common+com] * mb.datas[com*newC+col]
				}
				newDatas[col+row*newC] = sum
				ch <- ret{
					idx: col + row*newC,
					val: sum,
				}
			}(r, c, channel)
		}
	}

	for i := 0; i < newR*newC; i++ {
		item := <-channel
		newDatas[item.idx] = item.val
	}

	return NewMatrix(newR, newC, newDatas)
}

func (ma *Matrix) Mean(mb *Matrix) (*Matrix, error) {
	isSame := ma.sameShape(mb)
	if !isSame {
		return nil, errors.New("matrix different shape")
	}
	datas := []float64{}
	for i := range ma.datas {
		datas = append(datas, (ma.datas[i]+mb.datas[i])/2)
	}
	return NewMatrix(ma.rowNum, ma.colNum, datas)
}

func (ma *Matrix) RandMerge(mb *Matrix, orgRate float64) (*Matrix, error) {
	if orgRate < 0 || 1 < orgRate {
		return nil, errors.New("invalid rate range")
	}
	isSame := ma.sameShape(mb)
	if !isSame {
		return nil, errors.New("matrix different shape")
	}
	datas := []float64{}
	for i := range ma.datas {
		randf := rand.Float64()
		if randf < orgRate {
			datas = append(datas, ma.datas[i])
		} else {
			datas = append(datas, mb.datas[i])
		}
	}
	return NewMatrix(ma.rowNum, ma.colNum, datas)
}

func (m *Matrix) Func(f func(float64) (float64, error)) (*Matrix, error) {
	datas := []float64{}
	for i := range m.datas {
		val, err := f(m.datas[i])
		if err != nil {
			return nil, err
		}
		datas = append(datas, val)
	}
	return NewMatrix(m.rowNum, m.colNum, datas)
}

func (m *Matrix) Datas() []float64 {
	datas := []float64{}
	for i := range m.datas {
		datas = append(datas, m.datas[i])
	}
	return datas
}

func (ma *Matrix) sameShape(mb *Matrix) bool {
	return ma.colNum == mb.colNum && ma.rowNum == mb.rowNum
}
