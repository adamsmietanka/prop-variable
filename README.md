# Introduction
AWS Lambda function written in Go calculating efficiency of variable pitch propellers

# How does it work?
All the necessarry data is parsed from the query string of the GET method:
```
?max_speed=150&step_size=10&diameter=3.902&blades=3&cp=0.0902&prop_speed=20&power=800&angle=30&ratio=0.4
```

The program loads meshes from the ```\data``` directory based on the number of propeller blades. Then it generates an array of points to perform the interpolation.
The first step is to find the triangles each point lies on. The mesh is rectangular, so it makes things easier. 
Barycentric coordinates are then used to calculate the desired coordinate.