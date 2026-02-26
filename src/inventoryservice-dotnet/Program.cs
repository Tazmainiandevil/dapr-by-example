using Dapr;
using Dapr.Client;
using System.Text.Json;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddDaprClient();

var app = builder.Build();
app.UseCloudEvents();

app.MapPost("/orders", [Topic("pubsub", "orders")] async (Order incomingOrder) =>
{
    var order = incomingOrder ?? new Order("", 0);

    // Validation
    if (string.IsNullOrEmpty(order.OrderId))
    {
        Console.WriteLine("[INVENTORY] ❌ Received order with empty OrderId");
        return Results.BadRequest("OrderId is required");
    }
    if (order.Amount <= 0)
    {
        Console.WriteLine($"[INVENTORY] ❌ Received order {order.OrderId} with invalid amount: {order.Amount}");
        return Results.BadRequest("Amount must be positive");
    }

    Console.WriteLine($"[INVENTORY] 📦 Received order: {order.OrderId} (amount: ${order.Amount}) - {DateTime.Now:yyyy-MM-dd HH:mm:ss}");
    Console.WriteLine($"[INVENTORY] ✅ Order {order.OrderId} processed successfully");

    return Results.Ok();
});

// Add subscription endpoint for Dapr discovery
app.MapGet("/dapr/subscribe", () => 
{
    return new[] {
        new {
            pubsubname = "pubsub",
            topic = "orders",
            route = "/orders"
        }
    };
});

// Health check endpoint
app.MapGet("/healthz", () =>
{
    return Results.Ok(new { status = "healthy", service = "inventory" });
});

app.Run("http://0.0.0.0:8081");

public record Order(string OrderId, int Amount);
