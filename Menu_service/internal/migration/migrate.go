package migrations

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"menu/internal/model"
)

func Run(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	menuCol := db.Collection("menu")

	_, err := menuCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"category": 1},
	})
	if err != nil {
		log.Printf("⚠️ Failed to create index: %v", err)
	}
	count, _ := menuCol.CountDocuments(ctx, bson.M{})
	if count == 0 {
		items := []interface{}{
			model.MenuItem{
				Name:        "Classic Burger",
				Description: "Juicy beef patty with fresh lettuce, tomato, and our special sauce",
				Price:       9.99,
				Available:   true,
				Category:    "main-courses",
				ImageURL:    "https://images.unsplash.com/photo-1568901346375-23c9450c58cd?auto=format&fit=crop&w=1170&q=80",
			},
			model.MenuItem{
				Name:        "Margherita Pizza",
				Description: "Traditional Italian pizza with tomato sauce, mozzarella, and basil",
				Price:       12.99,
				Available:   true,
				Category:    "main-courses",
				ImageURL:    "https://images.unsplash.com/photo-1604068549290-dea0e4a305ca?auto=format&fit=crop&w=1074&q=80",
			},
			model.MenuItem{
				Name:        "Caesar Salad",
				Description: "Crisp romaine lettuce, croutons, and parmesan cheese with Caesar dressing",
				Price:       7.99,
				Available:   true,
				Category:    "appetizers",
				ImageURL:    "https://images.unsplash.com/photo-1550304943-4f24f54ddde9?auto=format&fit=crop&w=1170&q=80",
			},
			model.MenuItem{
				Name:        "Chicken Wings",
				Description: "Crispy chicken wings tossed in your choice of sauce",
				Price:       8.99,
				Available:   true,
				Category:    "appetizers",
				ImageURL:    "https://images.unsplash.com/photo-1567620832903-9fc6debc209f?auto=format&fit=crop&w=1080&q=80",
			},
			model.MenuItem{
				Name:        "Chocolate Lava Cake",
				Description: "Decadent chocolate cake with a gooey molten center",
				Price:       6.99,
				Available:   true,
				Category:    "desserts",
				ImageURL:    "https://images.unsplash.com/photo-1624353365286-3f8d62daad51?auto=format&fit=crop&w=1170&q=80",
			},
			model.MenuItem{
				Name:        "Iced Latte",
				Description: "Smooth espresso with cold milk over ice",
				Price:       3.99,
				Available:   true,
				Category:    "drinks",
				ImageURL:    "https://images.unsplash.com/photo-1517701550927-30cf4ba1dba5?auto=format&fit=crop&w=1170&q=80",
			},
			model.MenuItem{
				Name:        "Grilled Chicken Sandwich",
				Description: "Grilled chicken breast with lettuce and mayo",
				Price:       10.49,
				Available:   true,
				Category:    "main-courses",
				ImageURL:    "https://images.unsplash.com/photo-1597579018905-8c807adfbed4?w=500&auto=format&fit=crop&q=60",
			},
			model.MenuItem{
				Name:        "Vegetarian Wrap",
				Description: "Fresh vegetables wrapped in a soft tortilla",
				Price:       8.49,
				Available:   true,
				Category:    "main-courses",
				ImageURL:    "https://images.unsplash.com/photo-1592044903782-9836f74027c0?w=500&auto=format&fit=crop&q=60",
			},
			model.MenuItem{
				Name:        "Pepperoni Pizza",
				Description: "Classic pizza with spicy pepperoni and cheese",
				Price:       13.99,
				Available:   true,
				Category:    "main-courses",
				ImageURL:    "https://images.unsplash.com/photo-1628840042765-356cda07504e?w=500&auto=format&fit=crop&q=60",
			},
			model.MenuItem{
				Name:        "Garden Salad",
				Description: "Fresh garden vegetables with balsamic vinaigrette",
				Price:       6.99,
				Available:   true,
				Category:    "appetizers",
				ImageURL:    "https://images.unsplash.com/photo-1605291535126-2d71fea483c1?w=500&auto=format&fit=crop&q=60",
			},
			model.MenuItem{
				Name:        "Spaghetti Carbonara",
				Description: "Classic Italian pasta with creamy sauce",
				Price:       14.99,
				Available:   true,
				Category:    "main-courses",
				ImageURL:    "https://plus.unsplash.com/premium_photo-1674511582428-58ce834ce172?w=500&auto=format&fit=crop&q=60",
			},
			model.MenuItem{
				Name:        "Beef Tacos",
				Description: "Spiced beef with fresh toppings in a crispy shell",
				Price:       9.49,
				Available:   true,
				Category:    "main-courses",
				ImageURL:    "https://plus.unsplash.com/premium_photo-1661730314652-911662c0d86e?w=500&auto=format&fit=crop&q=60",
			},
			model.MenuItem{
				Name:        "Shrimp Cocktail",
				Description: "Chilled shrimp with tangy cocktail sauce",
				Price:       11.99,
				Available:   true,
				Category:    "appetizers",
				ImageURL:    "https://images.unsplash.com/photo-1691201659377-978b28daa417?w=500&auto=format&fit=crop&q=60",
			},
			model.MenuItem{
				Name:        "Tomato Soup",
				Description: "Rich and creamy tomato soup with croutons",
				Price:       5.49,
				Available:   true,
				Category:    "appetizers",
				ImageURL:    "https://images.unsplash.com/photo-1629978444632-9f63ba0eff47?w=500&auto=format&fit=crop&q=60",
			},
			model.MenuItem{
				Name:        "Berry Smoothie",
				Description: "Mixed berry smoothie with a touch of honey",
				Price:       4.99,
				Available:   true,
				Category:    "drinks",
				ImageURL:    "https://images.unsplash.com/photo-1553177595-4de2bb0842b9?w=500&auto=format&fit=crop&q=60",
			},
		}

		_, err := menuCol.InsertMany(ctx, items)
		if err != nil {
			log.Printf("Failed to insert seed data: %v", err)
		} else {
			log.Println("Seeded menu items")
		}
	}
}
