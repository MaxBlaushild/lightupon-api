from ortools.constraint_solver import pywrapcp
from math import radians, cos, sin, asin, sqrt

"""
________________________________________________________________________________________________________

This is meant to be the only public method on this module. It accepts a list of posts, and returns the 
same list sorted in a nice travelling salesman order. 

The caller does not have the flexibility the specify the "depot" (starting and ending point) of the new 
list. Instead, the depot is taken to be he first element in the provided list, and it should be the first 
element in the returned list. 

The bulk of the TSP code comes from here: https://developers.google.com/optimization/routing/tsp. 

In theory, a list containing any type can be passed to this method, under the constraint that each element 
has a numeric 'latitude' and 'longitude' property.
________________________________________________________________________________________________________
"""

def put_posts_into_tsp_order(posts):
  if len(posts) < 1:
    raise Exception('No posts provided.')
  dist_matrix = make_dist_matrix(posts)

  tsp_size = len(dist_matrix)
  num_routes = 1
  depot = 0

  # Create routing model
  if tsp_size > 0:
    routing = pywrapcp.RoutingModel(tsp_size, num_routes, depot)
    search_parameters = pywrapcp.RoutingModel.DefaultSearchParameters()
    # Create the distance callback.
    dist_callback = create_distance_callback(dist_matrix)
    routing.SetArcCostEvaluatorOfAllVehicles(dist_callback)
    # Solve the problem.
    assignment = routing.SolveWithParameters(search_parameters)
    if assignment:
      # Solution distance.
      print "Total distance: " + str(assignment.ObjectiveValue()) + " meters\n"
      # Display the solution.
      route_number = 0
      index = routing.Start(route_number) # Index of the variable for the starting node.
      new_posts = []

      while not routing.IsEnd(index):
        # Convert variable indices to node indices in the displayed route.
        index = assignment.Value(routing.NextVar(index))
        new_posts.append(posts[routing.IndexToNode(index)])

      return new_posts[-1:] + new_posts[:-1] # This is necessary because the TSP algo is putting the depot at the end of the quest.

    else:
      print 'No solution found.'
  else:
    raise Exception('Something went wrong while building the dist_matrix.')

def make_dist_matrix(posts):
  dist_matrix = []
  for post_a in posts:
    column = []
    for post_b in posts:
      if post_a == post_b:
        column.append(0)
      else:
        column.append(haversine_distance_in_integer_meters(post_a['longitude'], post_a['latitude'], post_b['longitude'], post_b['latitude']))
    dist_matrix.append(column)
  return dist_matrix

def haversine_distance_in_integer_meters(lon1, lat1, lon2, lat2):
    lon1, lat1, lon2, lat2 = map(radians, [lon1, lat1, lon2, lat2])

    dlon = lon2 - lon1 
    dlat = lat2 - lat1 
    a = sin(dlat/2)**2 + cos(lat1) * cos(lat2) * sin(dlon/2)**2
    c = 2 * asin(sqrt(a)) 
    r = 6371000 # Radius of earth in meters.
    return int(c * r)

def create_distance_callback(dist_matrix):

  def distance_callback(from_node, to_node):
    return int(dist_matrix[from_node][to_node])

  return distance_callback