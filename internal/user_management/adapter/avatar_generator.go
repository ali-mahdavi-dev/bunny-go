package adapter

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"math/rand"
	"os"
	"path/filepath"
	"sort"

	"github.com/disintegration/imaging"
)

type AvatarGenerator struct {
	avatarImagesPath string
	bodies           []string
	accessories      []string
	glasses          []string
	hats             []string
}

func NewAvatarGenerator() (*AvatarGenerator, error) {
	var err error
	g := &AvatarGenerator{}

	if g.bodies, err = loadLayerPaths("internal/user_management/assets/images/bodies"); err != nil {
		return nil, err
	}
	if g.accessories, err = loadLayerPaths("internal/user_management/assets/images/accessories"); err != nil {
		return nil, err
	}
	if g.glasses, err = loadLayerPaths("internal/user_management/assets/images/glasses"); err != nil {
		return nil, err
	}
	if g.hats, err = loadLayerPaths("internal/user_management/assets/images/hats"); err != nil {
		return nil, err
	}

	return g, nil
}

func (g *AvatarGenerator) Generate(identifier string) (image.Image, error) {
	if len(g.bodies) == 0 || len(g.accessories) == 0 || len(g.glasses) == 0 || len(g.hats) == 0 {
		return nil, errors.New("layers are empty")
	}

	seed := seedFromHash(identifier)
	r := rand.New(rand.NewSource(seed))

	layer0 := g.bodies[r.Intn(len(g.bodies))]
	layer1 := g.accessories[r.Intn(len(g.accessories))]
	layer2 := g.glasses[r.Intn(len(g.glasses))]
	layer3 := g.hats[r.Intn(len(g.hats))]

	img0, err := decodeImage(layer0)
	if err != nil {
		return nil, err
	}

	img1, err := decodeImage(layer1)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", filepath.Base(layer1), err)
	}
	img2, err := decodeImage(layer2)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", filepath.Base(layer2), err)
	}
	img3, err := decodeImage(layer3)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", filepath.Base(layer3), err)
	}

	avatar := imaging.Overlay(img0, img1, image.Point{0, 0}, 1.0)
	avatar = imaging.Overlay(avatar, img2, image.Point{0, 0}, 1.0)
	avatar = imaging.Overlay(avatar, img3, image.Point{0, 0}, 1.0)

	return avatar, nil
}

func loadLayerPaths(folder string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(folder, "*"))
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func decodeImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	return img, err
}

func seedFromHash(identifier string) int64 {
	h := md5.Sum([]byte(identifier))
	hexStr := hex.EncodeToString(h[:])
	var seed int64
	for i := 0; i < len(hexStr); i++ {
		seed = seed*16 + int64(hexStr[i])
	}
	return seed
}
