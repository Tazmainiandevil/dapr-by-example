using CommunityToolkit.Aspire.Hosting.Dapr;

var builder = DistributedApplication.CreateBuilder(args);

builder.AddProject<Projects.orderservice_dotnet>("ordersservice")
    .WithDaprSidecar(new DaprSidecarOptions
    {
        AppId = "orderservice",
        ResourcesPaths = [Path.Combine("../..", "components")]
    });

builder.AddProject<Projects.inventoryservice_dotnet>("inventoryservice")
    .WithDaprSidecar(new DaprSidecarOptions
    {
        AppId = "inventoryservice",
        ResourcesPaths = [Path.Combine("../..", "components")]
    });


builder.Build().Run();
