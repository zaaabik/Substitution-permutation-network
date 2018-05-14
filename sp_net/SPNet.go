package sp_net

import (
	"os"
	"math/rand"
	"math"
	"encoding/binary"
	"time"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"image/color"
	"gonum.org/v1/plot/vg/draw"
)

const blockSize = 32
const histbars = 256

type SPNet struct {
	rounds int
}

type Key struct {
	Bytes [32]byte
}

func (f SPNet) Encrypt(bytes []byte, key Key, rounds int, s [][]byte) ([]byte, error) {
	res := make([]byte, len(bytes))
	count := len(bytes)
	for i := 0; i < count; i += blockSize {
		tmp := [blockSize]byte{}

		if i+blockSize >= count {
			l := count - i
			zeros := make([]byte, l)
			bytes = append(bytes, zeros...)
			res = append(res, zeros...)
		}

		copy(tmp[:], bytes[i:i+blockSize])
		result := f.encryptBlock(tmp, key, s, rounds)
		copy(res[i:i+blockSize], result[:])
	}
	return res, nil
}

func (f SPNet) Decrypt(bytes []byte, key Key, rounds int, s [][]byte) ([]byte, error) {
	res := make([]byte, len(bytes))
	count := len(bytes)
	for i := 0; i < count; i += blockSize {
		tmp := [blockSize]byte{}

		if i+blockSize >= count {
			l := count - i
			zeros := make([]byte, l)
			bytes = append(bytes, zeros...)
			res = append(res, zeros...)
		}

		copy(tmp[:], bytes[i:i+blockSize])
		result := f.decryptBlock(tmp, key, s, rounds)
		copy(res[i:i+blockSize], result[:])
	}
	return res, nil
}

func (f SPNet) decryptBlock(block [blockSize]byte, key Key, s [][]byte, rounds int) (res [blockSize]byte) {
	res = [blockSize]byte{}

	for i := range block {
		res[i] = f.decryptByte(block[i], s, rounds)
		res[i] = res[i] ^ key.Bytes[i]
	}
	return res
}

func (f SPNet) encryptBlock(block [blockSize]byte, key Key, s [][]byte, rounds int) (res [blockSize]byte) {
	res = [blockSize]byte{}
	for i := range block {
		tmp := block[i] ^ key.Bytes[i]
		res[i] = f.encryptByte(tmp, s, rounds)
	}
	return res
}

func (f SPNet) encryptByte(block byte, s [][]byte, rounds int) (byte) {
	res := block
	size := len(s)
	for i := 0; i < rounds; i++ {
		lvl := i % size
		lvl2 := (i + 1) % size

		idx := find(s[lvl], res)
		res = s[lvl2][idx]
	}
	return res
}

func (f SPNet) decryptByte(block byte, s [][]byte, rounds int) (byte) {
	res := block
	size := len(s)
	for i := rounds - 1; i >= 0; i-- {
		lvl := (i + 1) % size
		lvl2 := (i) % size

		idx := find(s[lvl], res)
		res = s[lvl2][idx]
	}
	return res
}

func (f SPNet) GenerateBlock(path string, size int) {
	file, _ := os.Create(path)
	defer file.Close()

	rand.Seed(time.Now().UTC().UnixNano())

	if size == 1 {
		s := int(math.Pow(float64(2), float64(size*8)))
		blocks := make([]uint8, s)
		blocks2 := make([]uint8, s)
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
	} else if size == 2 {
		s := int(math.Pow(float64(2), float64(size*8)))
		blocks := make([]uint16, s)
		blocks2 := make([]uint16, s)
		p := rand.Perm(s)
		p2 := rand.Perm(s)

		for i := 0; i < s; i++ {
			blocks[i] = uint16(p[i])
			blocks2[i] = uint16(p2[i])
		}
		binary.Write(file, binary.LittleEndian, blocks)
		binary.Write(file, binary.LittleEndian, blocks2)
	} else if size == 3 {
		s := int(math.Pow(float64(2), float64(size*8)))
		blocks := make([]uint32, s)
		blocks2 := make([]uint32, s)
		p := rand.Perm(s)
		p2 := rand.Perm(s)

		for i := 0; i < s; i++ {
			blocks[i] = uint32(p[i])
			blocks2[i] = uint32(p2[i])
		}
		binary.Write(file, binary.LittleEndian, blocks)
		binary.Write(file, binary.LittleEndian, blocks2)
	}
}

func (f SPNet) ReadBlock1(path string) ([]uint8, []uint8, error) {
	file, _ := os.Open(path)
	defer file.Close()

	s := int(math.Pow(2, 8))

	buffer := make([]uint8, s)
	buffer2 := make([]uint8, s)
	file.Read(buffer)
	file.Read(buffer2)
	return buffer, buffer2, nil
}

func find(byte []byte, val byte) (int) {
	for k, v := range byte {
		if v == val {
			return k
		}
	}
	return -1
}

func cycleShift(x byte, k uint8) (byte) {
	return (x >> k) | (x << (8 - k))
}

func avg(byte []byte) float64 {
	avg := 0.0
	for _, v := range byte {
		avg += float64(v)
	}
	return avg / float64(len(byte))
}

func (f SPNet) AutoCorrelation(bytes []byte) []float64 {
	const size = 5
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
	for i, _ := range b {
		top += (float64(a[i]) - avgA) * (float64(b[i]) - avgB)
		bottom += math.Sqrt(math.Pow(float64(a[i])-avgA, 2) * math.Pow(float64(b[i])-avgB, 2))
	}
	return top / bottom
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
