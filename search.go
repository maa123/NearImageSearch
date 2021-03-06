package main

import(
	"image"
	"image/color"
	_"image/gif"
	_"image/jpeg"
	_"image/png"
	"archive/zip"
	"bufio"
	"fmt"
	"os"
	"math/bits"
	"sync"
	"time"
	"log"
	"strconv"

	"github.com/bamiaux/rez"
)

func LoadDB(filepath string) [][]uint64 {
	r, err := zip.OpenReader(filepath)
	defer r.Close()
	if err != nil {
		log.Fatal(err)
	}
	ldb := make([][]uint64, 12)
	for _, f := range r.File {
		i, _ := strconv.Atoi(f.Name)
		i--;
		ldb[i] = make([]uint64, 0, 34930)
		fp, _ := f.Open()
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			st, _ := strconv.ParseUint(scanner.Text(), 2, 64)
			ldb[i] = append(ldb[i], st)
		}
		fp.Close()
		if err = scanner.Err(); err != nil {
		}
	}
	return ldb
}

func rezResize(img image.Image) image.Image {
	drb := image.Rect(0, 0, 9, 8)
	var im image.Image
	switch t := img.(type) {
		case *image.YCbCr:
			im = image.NewYCbCr(drb, t.SubsampleRatio)
		case *image.Alpha:
			im = image.NewAlpha(drb)
		case *image.Alpha16:
			im = image.NewAlpha16(drb)
		case *image.Gray:
			im = image.NewGray(drb)
		case *image.Gray16:
			im = image.NewGray16(drb)
		case *image.NRGBA:
			im = image.NewNRGBA(drb)
		case *image.NRGBA64:
			im = image.NewNRGBA64(drb)
		case *image.RGBA:
			im = image.NewRGBA(drb)
		case *image.RGBA64:
			im = image.NewNRGBA64(drb)
		default:
			return nil
	}
	rez.Convert(im, img, rez.NewBicubicFilter())
	return im
}
func GetdHash(img image.Image) uint64 {
	var hash uint64
	var bit uint64 = 1
	var r, l uint32
	drb := image.Rect(0, 0, 9, 8)
	imgn := rezResize(img)
	gs := image.NewGray16(drb)
	for y := drb.Min.Y; y < drb.Max.Y; y++ {
		for x := drb.Min.X; x < drb.Max.X; x++ {
			c := color.Gray16Model.Convert(imgn.At(x, y))
			gray, _ := c.(color.Gray16)
			gs.Set(x, y, gray)
		}
	}
	for y := drb.Min.Y; y < drb.Max.Y; y++ {
		l, _, _, _ = gs.At(drb.Min.X, y).RGBA()
		for x := drb.Min.X + 1; x < drb.Max.X; x++ {
			r, _, _, _ = gs.At(x, y).RGBA()
			if r > l {
				hash |= bit
			}
			bit = bit << 1
			l = r
		}
	}
	return hash
}
func compare(a uint64, b uint64) uint8 {
	return uint8(bits.OnesCount64(a ^ b))
}

func Search(db [][]uint64, hash uint64, distanceValue uint8) (string, uint8) {
	reseps := make([]map[int]float64, 12)
	var wait sync.WaitGroup
	for ep := 0; ep < 12; ep++ {
		wait.Add(1)
		go func(epnm int, dist uint8) {
			defer wait.Done()
			reseps[epnm] = make(map[int]float64)
				for i, phash := range db[epnm] {
				if compare(phash, hash) < dist {
					reseps[epnm][i/24] = float64(i)
				}
			}
		}(ep, distanceValue)
	}
	wait.Wait()
	res := ""
	var lct uint8 = 0
	for ep := 0; ep < 12; ep++ {
		for _, frame := range reseps[ep] {
			if lct < 10 {
				res =  res + fmt.Sprintf("%d話の %.2f秒と一致しました\r\n", ep+1, (frame/24.0))
			}
			lct++
		}
	}
	if lct == 0 {
		return "一致するものはありませんでした", 0
	}
	return res, lct
}

func SearchFromPath(pt string, distanceValue uint8, hdb [][]uint64) (string, uint64) {
	file, err := os.Open(pt)
	if err != nil {
		return "Error", 0
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return "Image Error", 0
	}
	start := time.Now()
	hs := GetdHash(img)
	res, count := Search(hdb, hs, distanceValue)
	end := time.Now()
	if count == 0 {
		return fmt.Sprintf("一致するものはありませんでした\r\n処理時間:%f秒", (end.Sub(start)).Seconds()), hs
	}
	return fmt.Sprintf("%d件一致しました\r\n処理時間:%f秒\r\n\r\n", count, (end.Sub(start)).Seconds()) + res, hs
}

func SearchFromHash(hs uint64, distanceValue uint8, hdb [][]uint64) string {
	start := time.Now()
	res, count := Search(hdb, hs, distanceValue)
	end := time.Now()
	if count == 0 {
		return fmt.Sprintf("一致するものはありませんでした\r\n処理時間:%f秒", (end.Sub(start)).Seconds())
	}
	return fmt.Sprintf("%d件一致しました\r\n処理時間:%f秒\r\n\r\n", count, (end.Sub(start)).Seconds()) + res
}