# Build stage
FROM mcr.microsoft.com/dotnet/sdk:10.0 AS build
WORKDIR /src

# Copy project file
COPY ["inventoryservice-dotnet.csproj", "."]

# Restore dependencies
RUN dotnet restore "inventoryservice-dotnet.csproj"

# Copy source code
COPY . .

# Build the application
RUN dotnet build "inventoryservice-dotnet.csproj" -c Release -o /app/build

# Publish stage
FROM build AS publish
RUN dotnet publish "inventoryservice-dotnet.csproj" -c Release -o /app/publish

# Runtime stage
FROM mcr.microsoft.com/dotnet/aspnet:10.0
WORKDIR /app

# Copy published application from publish stage
COPY --from=publish /app/publish .

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
    CMD dotnet-health-check || exit 1

# Expose port
EXPOSE 8081

# Set environment
ENV ASPNETCORE_URLS=http://+:8081

# Run application
ENTRYPOINT ["dotnet", "inventoryservice-dotnet.dll"]
