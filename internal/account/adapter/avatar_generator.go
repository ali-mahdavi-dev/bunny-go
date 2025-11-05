package adapter

import (
	"bytes"
	"crypto/md5"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"io/fs"
	"math/rand"
	"path/filepath"
	"sort"

	"github.com/disintegration/imaging"

	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/infrastructure/logging"
)

type AvatarGenerator struct {
	bodies      []string
	accessories []string
	glasses     []string
	hats        []string
	imagesFS    embed.FS
}

func NewAvatarGenerator(imagesFS embed.FS) (*AvatarGenerator, error) {
	var err error
	g := &AvatarGenerator{imagesFS: imagesFS}

	logging.Debug("Loading avatar image layers").Log()

	if g.bodies, err = loadLayerPaths(imagesFS, "assets/images/bodies/*"); err != nil {
		logging.Error("Failed to load body images").
			WithError(err).
			Log()
		return nil, err
	}
	logging.Debug("Loaded body images").
		WithInt("count", len(g.bodies)).
		Log()

	if g.accessories, err = loadLayerPaths(imagesFS, "assets/images/accessories/*"); err != nil {
		logging.Error("Failed to load accessory images").
			WithError(err).
			Log()
		return nil, err
	}
	logging.Debug("Loaded accessory images").
		WithInt("count", len(g.accessories)).
		Log()

	if g.glasses, err = loadLayerPaths(imagesFS, "assets/images/glasses/*"); err != nil {
		logging.Error("Failed to load glasses images").
			WithError(err).
			Log()
		return nil, err
	}
	logging.Debug("Loaded glasses images").
		WithInt("count", len(g.glasses)).
		Log()

	if g.hats, err = loadLayerPaths(imagesFS, "assets/images/hats/*"); err != nil {
		logging.Error("Failed to load hat images").
			WithError(err).
			Log()
		return nil, err
	}
	logging.Debug("Loaded hat images").
		WithInt("count", len(g.hats)).
		Log()

	logging.Info("Avatar generator initialized successfully").Log()
	return g, nil
}

func (g *AvatarGenerator) Generate(identifier string) (image.Image, error) {
	if len(g.bodies) == 0 || len(g.accessories) == 0 || len(g.glasses) == 0 || len(g.hats) == 0 {
		err := errors.New("layers are empty")
		logging.Error("Avatar generation failed: layers are empty").
			WithString("identifier", identifier).
			WithInt("bodies_count", len(g.bodies)).
			WithInt("accessories_count", len(g.accessories)).
			WithInt("glasses_count", len(g.glasses)).
			WithInt("hats_count", len(g.hats)).
			WithError(err).
			Log()
		return nil, err
	}

	seed := seedFromHash(identifier)
	r := rand.New(rand.NewSource(seed))

	layer0 := g.bodies[r.Intn(len(g.bodies))]
	layer1 := g.accessories[r.Intn(len(g.accessories))]
	layer2 := g.glasses[r.Intn(len(g.glasses))]
	layer3 := g.hats[r.Intn(len(g.hats))]

	img0, err := decodeImage(g.imagesFS, layer0)
	if err != nil {
		logging.Error("Failed to decode body image").
			WithString("identifier", identifier).
			WithString("file", filepath.Base(layer0)).
			WithError(err).
			Log()
		return nil, fmt.Errorf("failed to decode %s: %w", filepath.Base(layer0), err)
	}

	img1, err := decodeImage(g.imagesFS, layer1)
	if err != nil {
		logging.Error("Failed to decode accessory image").
			WithString("identifier", identifier).
			WithString("file", filepath.Base(layer1)).
			WithError(err).
			Log()
		return nil, fmt.Errorf("failed to decode %s: %w", filepath.Base(layer1), err)
	}
	img2, err := decodeImage(g.imagesFS, layer2)
	if err != nil {
		logging.Error("Failed to decode glasses image").
			WithString("identifier", identifier).
			WithString("file", filepath.Base(layer2)).
			WithError(err).
			Log()
		return nil, fmt.Errorf("failed to decode %s: %w", filepath.Base(layer2), err)
	}

	img3, err := decodeImage(g.imagesFS, layer3)
	if err != nil {
		logging.Error("Failed to decode hat image").
			WithString("identifier", identifier).
			WithString("file", filepath.Base(layer3)).
			WithError(err).
			Log()
		return nil, fmt.Errorf("failed to decode %s: %w", filepath.Base(layer3), err)
	}

	avatar := imaging.Overlay(img0, img1, image.Point{0, 0}, 1.0)
	avatar = imaging.Overlay(avatar, img2, image.Point{0, 0}, 1.0)
	avatar = imaging.Overlay(avatar, img3, image.Point{0, 0}, 1.0)

	logging.Debug("Avatar generated successfully").
		WithString("identifier", identifier).
		Log()

	return avatar, nil
}

func loadLayerPaths(folder embed.FS, path string) ([]string, error) {
	files, err := fs.Glob(folder, path)
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func decodeImage(folder embed.FS, path string) (image.Image, error) {
	file, err := fs.ReadFile(folder, path)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(file))
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
