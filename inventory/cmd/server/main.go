package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	inventoryv1 "github.com/clava1096/rocket-service/shared/pkg/proto/inventory/v1"
)

const grpcPort = ":50051"

type inventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer
	mu sync.Mutex

	inventory map[string]*inventoryv1.Part
}

func (s *inventoryService) GetPart(_ context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	part, ok := s.inventory[req.Uuid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "part: %s not found", req.Uuid)
	}
	return &inventoryv1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *inventoryService) ListParts(_ context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var parts []*inventoryv1.Part

	for _, part := range s.inventory {
		parts = append(parts, part)
	}

	filter := req.GetFilter()
	if filter == nil {
		return &inventoryv1.ListPartsResponse{Parts: parts}, nil
	}

	if uuids := filter.GetUuids(); len(uuids) > 0 {
		parts = filterByUUID(parts, uuids)
	}

	if names := filter.GetNames(); len(names) > 0 {
		parts = filterByNames(parts, names)
	}

	if category := filter.GetCategories(); len(category) > 0 {
		parts = filterByCategories(parts, category)
	}

	if countries := filter.GetManufacturerCountries(); len(countries) > 0 {
		parts = filterByCountries(parts, countries)
	}

	if tags := filter.GetTags(); len(tags) > 0 {
		parts = FilterByTags(parts, tags)
	}

	return &inventoryv1.ListPartsResponse{
		Parts: parts,
	}, nil
}

func filterByCategories(parts []*inventoryv1.Part, categories []inventoryv1.Category) []*inventoryv1.Part {
	n := make(map[inventoryv1.Category]struct{}, len(categories))
	for _, cat := range categories {
		n[cat] = struct{}{}
	}
	var filtered []*inventoryv1.Part
	for _, part := range parts {
		if _, found := n[part.GetCategory()]; found {
			filtered = append(filtered, part)
		}
	}
	return filtered
}

func FilterByTags(parts []*inventoryv1.Part, tags []string) []*inventoryv1.Part {
	n := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		n[tag] = struct{}{}
	}
	var filtered []*inventoryv1.Part
	for _, part := range parts {
		for _, tag := range part.GetTags() {
			if _, found := n[tag]; found {
				filtered = append(filtered, part)
				break
			}
		}
	}
	return filtered
}

func filterByCountries(parts []*inventoryv1.Part, countries []string) []*inventoryv1.Part {
	n := make(map[string]struct{}, len(countries))
	for _, country := range countries {
		n[country] = struct{}{}
	}
	var filtered []*inventoryv1.Part
	for _, part := range parts {
		if man := part.GetManufacturer(); man != nil {
			if _, found := n[man.GetCountry()]; found {
				filtered = append(filtered, part)
			}
		}
	}
	return filtered
}

func filterByNames(parts []*inventoryv1.Part, names []string) []*inventoryv1.Part {
	n := make(map[string]struct{}, len(names))
	for _, name := range names {
		n[name] = struct{}{}
	}
	var filtered []*inventoryv1.Part
	for _, part := range parts {
		if _, found := n[part.GetName()]; found {
			filtered = append(filtered, part)
		}
	}
	return filtered
}

func filterByUUID(parts []*inventoryv1.Part, uuids []string) []*inventoryv1.Part {
	uuid := make(map[string]struct{}, len(uuids))
	for _, p := range uuids {
		uuid[p] = struct{}{}
	}
	var filtered []*inventoryv1.Part
	for _, part := range parts {
		if _, found := uuid[part.GetUuid()]; found {
			filtered = append(filtered, part)
		}
	}
	return filtered
}

func main() {
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v", cerr)
		}
	}()

	s := grpc.NewServer()
	service := &inventoryService{
		inventory: make(map[string]*inventoryv1.Part),
	}
	//Двигатель
	service.inventory["123e4567-e89b-12d3-a456-426614174000"] = &inventoryv1.Part{
		Uuid:        "123e4567-e89b-12d3-a456-426614174000",
		Name:        "Main Engine",
		Description: "High-thrust rocket engine",
		Price:       150000.0,
		Category:    inventoryv1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryv1.Dimensions{
			Length: 200.0,
			Width:  100.0,
			Height: 150.0,
			Weight: 500.0,
		},
		Manufacturer: &inventoryv1.Manufacturer{
			Name:    "CosmoEngines Inc.",
			Country: "USA",
			Website: "https://cosmoengines.example.com",
		},
		Tags: []string{"engine", "thrust", "rocket"},
	}
	//Крыло
	service.inventory["123e4567-e89b-12d3-a456-426614174001"] = &inventoryv1.Part{
		Uuid:        "123e4567-e89b-12d3-a456-426614174001",
		Name:        "Aerodynamic Wing",
		Description: "Lightweight composite wing for orbital stability",
		Price:       42000.50,
		Category:    inventoryv1.Category_CATEGORY_WING,
		Dimensions: &inventoryv1.Dimensions{
			Length: 300.0,
			Width:  80.0,
			Height: 20.0,
			Weight: 120.0,
		},
		Manufacturer: &inventoryv1.Manufacturer{
			Name:    "AeroStructures Ltd.",
			Country: "Germany",
			Website: "https://aerostructures.example.com",
		},
		Tags: []string{"wing", "aero", "composite"},
	}

	// Иллюминатор
	service.inventory["123e4567-e89b-12d3-a456-426614174002"] = &inventoryv1.Part{
		Uuid:        "123e4567-e89b-12d3-a456-426614174002",
		Name:        "Reinforced Porthole",
		Description: "Triple-glazed viewport for deep-space observation",
		Price:       8500.75,
		Category:    inventoryv1.Category_CATEGORY_PORTHOLE,
		Dimensions: &inventoryv1.Dimensions{
			Length: 50.0,
			Width:  50.0,
			Height: 10.0,
			Weight: 25.0,
		},
		Manufacturer: &inventoryv1.Manufacturer{
			Name:    "ViewSpace Optics",
			Country: "France",
			Website: "https://viewspace.example.com",
		},
		Tags: []string{"window", "observation", "safe"},
	}
	inventoryv1.RegisterInventoryServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("Starting gRPC server at %s", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	s.GracefulStop()
	log.Println("Server gracefully stopped")
}
