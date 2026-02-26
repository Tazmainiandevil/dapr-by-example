using System.Text;
using Dapr.Client;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddDaprClient();

var app = builder.Build();

app.MapPost("/orders", async (Order order, DaprClient dapr) =>
{
    // Validation
    if (string.IsNullOrEmpty(order.OrderId))
    {
        return Results.BadRequest("OrderId is required");
    }
    if (order.Amount <= 0)
    {
        return Results.BadRequest("Amount must be positive");
    }

    try
    {
        await dapr.SaveStateAsync("statestore", order.OrderId, order);
        await dapr.PublishEventAsync("pubsub", "orders", order);

        var metadata = new Dictionary<string, string>
        {
            ["blobName"] = $"{order.OrderId}.txt",
            ["key"] = $"{order.OrderId}.txt",
            ["fileName"] = $"{order.OrderId}.txt"
        };

        await dapr.InvokeBindingAsync(
            "storage",
            "create",
            Encoding.UTF8.GetBytes($"Order receipt for {order.OrderId}"),
            metadata
        );

        Console.WriteLine($"Order {order.OrderId} created successfully");
        return Results.Accepted();
    }
    catch (Exception ex)
    {
        Console.WriteLine($"Error processing order {order.OrderId}: {ex.Message}");
        return Results.Problem("Failed to process order");
    }
});

app.MapGet("/orders/{orderId}", async (string orderId, DaprClient dapr) =>
{
    var order = await dapr.GetStateAsync<Order>("statestore", orderId);
    if (order == null)
    {
        return Results.NotFound();
    }
    return Results.Ok(order);
});

// Add subscription endpoint for Dapr discovery (returns empty array since this is a publisher only)
app.MapGet("/dapr/subscribe", () => 
{
    return new object[0]; // Empty array - no subscriptions
});

// Health check endpoint
app.MapGet("/healthz", async (DaprClient dapr) =>
{
    try
    {
        var metadata = await dapr.GetMetadataAsync();
        return Results.Ok(new { status = "healthy", dapr = metadata });
    }
    catch
    {
        return Results.StatusCode(503);
    }
});

app.Run("http://0.0.0.0:8080");

public record Order(string OrderId, int Amount);