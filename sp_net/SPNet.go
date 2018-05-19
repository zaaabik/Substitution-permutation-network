package sp_net

import (
	"os"
	"math/rand"
	"math"
	"time"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"image/color"
	"gonum.org/v1/plot/vg/draw"
	"reflect"
)

const BlockSize = 32
const histbars = 256

type SPNet struct {
}

func (f SPNet) Encrypt(bytes []byte, p [][]byte, rounds int, s [][][]byte) ([]byte, error) {
	res := make([]byte, len(bytes))
	count := len(bytes)
	for i := 0; i < count; i += BlockSize {
		tmp := [BlockSize]byte{}
		if i+BlockSize >= count {
			l := count - i
			zeros := make([]byte, l)
			bytes = append(bytes, zeros...)
			res = append(res, zeros...)
		}

		copy(tmp[:], bytes[i:i+BlockSize])
		result := f.encryptBlock(tmp, p, s, rounds)
		copy(res[i:i+BlockSize], result[:])
	}
	return res, nil
}

func (f SPNet) Decrypt(bytes []byte, p [][]byte, rounds int, s [][][]byte) ([]byte, error) {
	res := make([]byte, len(bytes))
	count := len(bytes)
	for i := 0; i < count; i += BlockSize {
		tmp := [BlockSize]byte{}

		if i+BlockSize >= count {
			l := count - i
			zeros := make([]byte, l)
			bytes = append(bytes, zeros...)
			res = append(res, zeros...)
		}
		copy(tmp[:], bytes[i:i+BlockSize])
		result := f.decryptBlock(tmp, p, s, rounds)
		copy(res[i:i+BlockSize], result[:])
	}
	return res, nil
}

func (f SPNet) decryptBlock(block [BlockSize]byte, p [][]byte, s [][][]byte, rounds int) (res [BlockSize]byte) {
	res = block
	for i := rounds - 1; i >= 0; i-- {
		resAfterPermutation := [BlockSize]byte{}
		key := InveseKey(p[i])
		for j := range block {
			resAfterPermutation[j] = res[key[j]]
		}
		res = resAfterPermutation
		for c := range block {
			res[c] = f.decryptByte(res[c], s[i][c], s[i+1][c])
		}
	}

	return
}

func InveseKey(p []byte) ([]byte){
	key:= make([]byte, len(p))
	for i:= 0; i < len(p); i++ {
		key[ p[i]] = byte(i)
	}
	return key
}

func (f SPNet) encryptBlock(block [BlockSize]byte, p [][]byte, s [][][]byte, rounds int) (res [BlockSize]byte) {
	res = block
	for i := 0; i < rounds; i++ {
		for c := range block {
			res[c] = f.encryptByte(block[c], s[i][c], s[i+1][c])
		}
		resAfterPermutation := [BlockSize]byte{}
		for j := range block {
			resAfterPermutation[j] = res[p[i][j]]
		}

		res = resAfterPermutation
	}
	return
}

func (f SPNet) encryptByte(block byte, s []byte, s2 []byte) (byte) {
	res := block
	idx := find(s, res)
	res = s2[idx]
	return res
}

func (f SPNet) decryptByte(block byte, s []byte, s2 []byte) (byte) {
	res := block
	idx := find(s2, res)
	res = s[idx]
	return res
}

func (f SPNet) GenerateSBlock(path string, count int) {
	file, _ := os.Create(path)
	defer file.Close()

	rand.Seed(time.Now().UTC().UnixNano())

	s := int(math.Pow(float64(2), float64(8)))
	blocks := make([]uint8, s)
	blocks2 := make([]uint8, s)
	for c := 0; c < count*BlockSize; c++ {
		p := rand.Perm(s)
		p2 := rand.Perm(s)
		for i := 0; i < s; i++ {
			blocks[i] = uint8(p[i])
		}

		for i := 0; i < s; i++ {
			blocks2[i] = uint8(p2[i])
		}
		file.Write(blocks)
		file.Write(blocks2)
	}
}

func (f SPNet) ReadSBlocks(path string) ([][][]byte, error) {
	file, _ := os.Open(path)
	defer file.Close()
	stats, _ := file.Stat()
	count := stats.Size() / BlockSize / 256
	blocks := make([][][]byte, count)
	for j := 0; j < int(count); j++ {
		blocks[j] = make([][]byte, BlockSize)
	}
	for i := 0; i < int(count); i += 2 {
		for j := 0; j < BlockSize; j++ {
			s1 := make([]byte, 256)
			s2 := make([]byte, 256)
			file.Read(s1)
			file.Read(s2)
			blocks[i][j] = s1
			blocks[i+1][j] = s2
		}
	}

	return blocks, nil
}

func find(byte []byte, val byte) (int) {
	for k, v := range byte {
		if v == val {
			return k
		}
	}
	return -1
}

func rightCycleShift(x byte, k uint8) (byte) {
	return (x >> k) | (x << (8 - k))
}

func leftCycleShift(x byte, k uint8) (byte) {
	return (x << k) | (x >> (8 - k))
}

func avg(byte []byte) float64 {
	avg := 0.0
	for _, v := range byte {
		avg += float64(v)
	}
	return avg / float64(len(byte))
}

func (f SPNet) AutoCorrelation(bytes []byte, size int) []float64 {
	result := make([]float64, size)
	for i := 0; i < size; i++ {
		a := bytes[i:]
		b := bytes[:len(bytes)-i]
		result[i] = f.Correlation(a, b)
	}
	return result
}

func (f SPNet) Correlation(a, b []byte) float64 {
	top := 0.
	bottom := 0.
	avgA := avg(a)
	avgB := avg(b)
	l := math.Min(float64(len(a)), float64(len(b)))
	for i := 0; i < int(l); i++ {
		top += (float64(a[i]) - avgA) * (float64(b[i]) - avgB)
		bottom += math.Sqrt(math.Pow(float64(a[i])-avgA, 2) * math.Pow(float64(b[i])-avgB, 2))
	}
	return top / bottom
}

func (f SPNet) GeneratePBlocks(count int, path string) (error) {
	rand.Seed(time.Now().UnixNano())
	buffer := make([]byte, 0)
	for c := 0; c < count; c++ {
		list := rand.Perm(BlockSize)
		tmp := make([]byte, BlockSize)
		for i := range list {
			tmp[i] = byte(list[i])
		}
		buffer = append(buffer, tmp...)
	}
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}
	file.Write(buffer)
	return nil
}
func Test(bytes []byte, path string) {
	x := bytes[1:]
	y := bytes[:len(bytes)-1]
	scatterData := getPoint(x, y)
	minZ, maxZ := math.Inf(1), math.Inf(-1)
	for _, xyz := range scatterData {
		if xyz.Z > maxZ {
			maxZ = xyz.Z
		}
		if xyz.Z < minZ {
			minZ = xyz.Z
		}
	}
	// Create a new plot, set its title and
	// axis labels.
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Points Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	// Draw a grid behind the data
	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(scatterData)
	s.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		c := color.RGBA{R: 196, B: 128, A: 255}
		return draw.GlyphStyle{Color: c, Radius: 1, Shape: draw.CircleGlyph{}}
	}
	if err != nil {
		panic(err)
	}
	// Add the plotters to the plot, with a legend
	// entry for each
	p.Add(s)
	p.Legend.Add("scatter", s)

	// Save the plot to a PNG file.
	if err := p.Save(20*vg.Inch, 20*vg.Inch, path); err != nil {
		panic(err)
	}
}

// randomPoints returns some random x, y points.
func getPoint(x, y []byte) plotter.XYZs {
	pts := make(plotter.XYZs, len(x))
	for i := range x {
		pts[i].X = float64(x[i])
		pts[i].Y = float64(y[i])
		pts[i].Z = 0.1

	}
	return pts
}

func MathExpected(byte []byte) float64 {
	res := avg(byte)
	return res
}

func Dispersion(bytes []byte) float64 {
	average := avg(bytes)
	res := make([]float64, len(bytes))
	for i := range bytes {
		res[i] = math.Pow(float64(bytes[i])-average, 2)
	}

	result := 0.0
	for i := range res {
		result += res[i]
	}
	return result / float64(len(res))
}

func MakeHist(path string, bytes []byte) {
	rand.Seed(int64(0))
	v := make(plotter.Values, len(bytes))
	for i := range v {
		v[i] = float64(bytes[i])
	}

	// Make a plot and set its title.
	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Histogram"
	h, err := plotter.NewHist(v, histbars)
	if err != nil {
		panic(err)
	}

	h.Normalize(1)
	p.Add(h)

	// Save the plot to a PNG file.
	if err := p.Save(20*vg.Inch, 20*vg.Inch, path); err != nil {
		panic(err)
	}
}

func Shuffle(slice interface{}) {
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()
	for i := length - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		swap(i, j)
	}
}
