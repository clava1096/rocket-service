//go:build integration

package integration

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	inventoryv1 "github.com/clava1096/rocket-service/shared/pkg/proto/inventory/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ = Describe("InventoryService E2E", func() {
	var (
		ctx       context.Context
		cancel    context.CancelFunc
		inventory inventoryv1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()))

		Expect(err).NotTo(HaveOccurred(), "Except success connect to gRPC server")

		inventory = inventoryv1.NewInventoryServiceClient(conn)
	})

	AfterEach(func() {

		cancel()
	})

	Describe("GetPart", func() {
		Context("when the part exists in MongoDB", func() {
			var partUUID string
			var expectedPart *inventoryv1.Part

			BeforeEach(func() {
				// Seed data directly into MongoDB using test helper
				expectedPart = env.GetTestPartInfo()
				var err error
				partUUID, err = env.InsertTestDetailsWithData(ctx, expectedPart)
				Expect(err).NotTo(HaveOccurred(), "expected successful insertion of test part")
			})

			It("should return the part by UUID via gRPC", func() {
				resp, err := inventory.GetPart(ctx, &inventoryv1.GetPartRequest{
					Uuid: partUUID,
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.GetPart()).NotTo(BeNil())
				Expect(resp.GetPart().Uuid).To(Equal(partUUID))
				Expect(resp.GetPart().Name).To(Equal(expectedPart.Name))
				Expect(resp.GetPart().Description).To(Equal(expectedPart.Description))
				Expect(resp.GetPart().Price).To(Equal(expectedPart.Price))
				Expect(resp.GetPart().Category).To(Equal(expectedPart.Category))
				Expect(resp.GetPart().CreatedAt).NotTo(BeNil())
			})

			It("should return part with correct nested structures", func() {
				resp, err := inventory.GetPart(ctx, &inventoryv1.GetPartRequest{
					Uuid: partUUID,
				})

				Expect(err).NotTo(HaveOccurred())
				part := resp.GetPart()

				// Verify dimensions
				Expect(part.Dimensions).NotTo(BeNil())
				Expect(part.Dimensions.Length).To(Equal(expectedPart.Dimensions.Length))
				Expect(part.Dimensions.Weight).To(Equal(expectedPart.Dimensions.Weight))

				// Verify manufacturer
				Expect(part.Manufacturer).NotTo(BeNil())
				Expect(part.Manufacturer.Name).To(Equal(expectedPart.Manufacturer.Name))
				Expect(part.Manufacturer.Country).To(Equal(expectedPart.Manufacturer.Country))

				// Verify tags
				Expect(part.Tags).To(ConsistOf(expectedPart.Tags))
			})
		})

		Context("when the part does not exist", func() {
			It("should return an error for non-existent UUID", func() {
				nonExistentUUID := "00000000-0000-0000-0000-000000000000"

				resp, err := inventory.GetPart(ctx, &inventoryv1.GetPartRequest{
					Uuid: nonExistentUUID,
				})

				// Adjust expectation based on your service error handling
				if err != nil {
					Expect(err.Error()).To(ContainSubstring("not found"))
				} else {
					Expect(resp.GetPart()).To(BeNil())
				}
			})
		})
	})

	Describe("ListParts", func() {
		BeforeEach(func() {
			// Seed multiple parts for filtering/pagination tests
			for i := 0; i < 5; i++ {
				part := env.GetTestPartInfo()
				if i < 3 {
					part.Category = inventoryv1.Category_CATEGORY_PORTHOLE
				} else {
					part.Category = inventoryv1.Category_CATEGORY_ENGINE
				}
				_, err := env.InsertTestDetailsWithData(ctx, part)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should return all parts when no filters applied", func() {
			resp, err := inventory.ListParts(ctx, &inventoryv1.ListPartsRequest{
				Filter: &inventoryv1.PartsFilter{},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.GetParts()).To(HaveLen(5))
		})

		It("should filter parts by category", func() {
			resp, err := inventory.ListParts(ctx, &inventoryv1.ListPartsRequest{
				Filter: &inventoryv1.PartsFilter{
					Categories: []inventoryv1.Category{
						inventoryv1.Category_CATEGORY_PORTHOLE,
					},
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.GetParts()).To(HaveLen(3))
			for _, part := range resp.GetParts() {
				Expect(part.Category).To(Equal(inventoryv1.Category_CATEGORY_PORTHOLE))
			}
		})

		It("should filter parts by multiple categories", func() {
			resp, err := inventory.ListParts(ctx, &inventoryv1.ListPartsRequest{
				Filter: &inventoryv1.PartsFilter{
					// Фильтр по нескольким категориям
					Categories: []inventoryv1.Category{
						inventoryv1.Category_CATEGORY_PORTHOLE,
						inventoryv1.Category_CATEGORY_ENGINE,
					},
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.GetParts()).To(HaveLen(5))
		})

		It("should filter parts by tags", func() {
			// Вставляем часть с известным тегом
			part := env.GetTestPartInfo()
			part.Tags = []string{"special-tag", "another-tag"}
			part.Category = inventoryv1.Category_CATEGORY_WING
			_, err := env.InsertTestDetailsWithData(ctx, part)
			Expect(err).NotTo(HaveOccurred())

			resp, err := inventory.ListParts(ctx, &inventoryv1.ListPartsRequest{
				Filter: &inventoryv1.PartsFilter{
					// Фильтр по тегам
					Tags: []string{"special-tag"},
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.GetParts()).NotTo(BeEmpty())
			// Проверяем, что хотя бы одна запись содержит нужный тег
			found := false
			for _, p := range resp.GetParts() {
				for _, tag := range p.Tags {
					if tag == "special-tag" {
						found = true
						break
					}
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should filter parts by manufacturer country", func() {
			// Вставляем часть с известной страной
			part := env.GetTestPartInfo()
			part.Manufacturer.Country = "Mars"
			part.Category = inventoryv1.Category_CATEGORY_FUEL
			_, err := env.InsertTestDetailsWithData(ctx, part)
			Expect(err).NotTo(HaveOccurred())

			resp, err := inventory.ListParts(ctx, &inventoryv1.ListPartsRequest{
				Filter: &inventoryv1.PartsFilter{
					// Фильтр по стране производителя
					ManufacturerCountries: []string{"Mars"},
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.GetParts()).NotTo(BeEmpty())
			for _, p := range resp.GetParts() {
				Expect(p.Manufacturer.Country).To(Equal("Mars"))
			}
		})

		It("should filter parts by name (partial match if implemented)", func() {
			// Вставляем часть с уникальным именем
			uniqueName := "SuperWidget-" + gofakeit.UUID()
			part := env.GetTestPartInfo()
			part.Name = uniqueName
			_, err := env.InsertTestDetailsWithData(ctx, part)
			Expect(err).NotTo(HaveOccurred())

			resp, err := inventory.ListParts(ctx, &inventoryv1.ListPartsRequest{
				Filter: &inventoryv1.PartsFilter{
					// Фильтр по имени (зависит от реализации: точное совпадение или LIKE)
					Names: []string{uniqueName},
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.GetParts()).NotTo(BeEmpty())
			Expect(resp.GetParts()[0].Name).To(Equal(uniqueName))
		})

		It("should filter parts by UUIDs list", func() {
			// Получаем UUID двух случайных частей
			allResp, err := inventory.ListParts(ctx, &inventoryv1.ListPartsRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(allResp.GetParts()).To(HaveLen(5))

			targetUUIDs := []string{
				allResp.GetParts()[0].Uuid,
				allResp.GetParts()[2].Uuid,
			}

			resp, err := inventory.ListParts(ctx, &inventoryv1.ListPartsRequest{
				Filter: &inventoryv1.PartsFilter{
					// Фильтр по списку UUID
					Uuids: targetUUIDs,
				},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(resp.GetParts()).To(HaveLen(2))

			// Проверяем, что вернулись именно те UUID
			returnedUUIDs := make([]string, len(resp.GetParts()))
			for i, p := range resp.GetParts() {
				returnedUUIDs[i] = p.Uuid
			}
			Expect(returnedUUIDs).To(ConsistOf(targetUUIDs))
		})
	})

	Describe("GetPart with metadata", func() {
		It("should return part with metadata fields correctly mapped", func() {
			partInfo := env.GetTestPartInfo()
			partInfo.Metadata = map[string]*inventoryv1.Value{
				"color": {
					Kind: &inventoryv1.Value_StringValue{StringValue: "blue"},
				},
				"weight_kg": {
					Kind: &inventoryv1.Value_DoubleValue{DoubleValue: 12.5},
				},
			}

			uuid, err := env.InsertTestDetailsWithData(ctx, partInfo)
			Expect(err).NotTo(HaveOccurred())

			resp, err := inventory.GetPart(ctx, &inventoryv1.GetPartRequest{Uuid: uuid})
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.GetPart().Metadata).To(HaveKey("color"))
			Expect(resp.GetPart().Metadata["color"].GetStringValue()).To(Equal("blue"))
			Expect(resp.GetPart().Metadata["weight_kg"].GetDoubleValue()).To(Equal(12.5))
		})
	})
})
