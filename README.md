# Measurements API

This API, created for Check24 interview preparation, manages measurements linked to experiments and sensors, supporting basic CRUD operations and time range queries.

## Endpoints

- **GET /measurements**: Retrieve all measurements.
- **POST /measurements**: Create a new measurement.
- **GET /measurements/:id**: Retrieve a measurement by ID.
- **PUT /measurements/:id**: Update a measurement by ID.
- **DELETE /measurements/:id**: Delete a measurement by ID.
- **GET /measurements/minmax**: Get the min and max measurement values.
- **GET /experiments/:exp/measurements**: Get measurements for a specific experiment (optionally with a time range).
  - Example: `GET /experiments/Exp2/measurements?startTime=2025-02-24+00:00:00&endTime=2025-02-24+11:26:00`
  
## Database Schema
- **Experiments**: Stores experiment details.
- **Sensors**: Stores sensor data.
- **Measurements**: Stores measurement data linked to experiments and sensors.
