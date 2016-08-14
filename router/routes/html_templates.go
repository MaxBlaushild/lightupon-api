package routes

const guess_what_i_said = `
<html>
  artisinal honey loaves. that's all we sell.
</html>
`

const trips_list_template = `
<style>
  .bold {
    font-weight: bold;
  }

  .block_container {
    border: 1px solid black;
    padding-left:3px;
  }

  .all-trips {
    width: 400px;
    float: left;

  }

  .scenes_for_trip {
    width: 500px;
    border: 1px solid black;
    padding-left:5px;

  }

  .add_trip {
    margin-left: 420px;
    width:300px;
  }

  .button {
    border: 1px solid black;
    background-color:#EDA1E7;
    width:50px;
  }

</style>
<html>
  <body>
    <div style="width: 100%; overflow: hidden;">
      <div class="all-trips block_container">
        <p class="bold"> YOUR TRIPS<p>
          {{range $index, $element := .}}
            <p class="bold"><a href="/lightupon/admin/trips/{{$element.ID}}">{{$element.Title}}</a></p>
            <p>{{$element.Description}}</p>
            <p><img src="{{$element.ImageUrl}}" height="50" width="150"/></p>
            <p><span class="delete_trip button" id="trip_{{$element.ID}}">delete trip</span></p>
          {{end}}
      </div>
      <div class="add_trip block_container" >
        <p class="bold"> CREATE TRIP </p>
        <p>Trip Title: <input type="text" id="input-trip_title"/></p>
        <p>Trip Description: <input type="text" id="input-trip_description"/></p>
        <p>Trip Image URL: <input type="text" id="input-trip_image_url"/></p>
        <p class="submit_trip button">Submit</p>
      </div>
    </div>
  </body>
</html>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>
<script>
(function(){

$('.submit_trip').on('click', function(){
  post_trip();
  window.location.reload(false);
})

$('.delete_trip').each(function(index, element){
  var delete_trip_element = $(element);
  delete_trip_element.on('click', function(element_1){
    var trip_id = delete_trip_element.attr('id').split('_')[1];
    console.log(trip_id)
      
    $.ajax({
      method: "DELETE",
      url:"/lightupon/admin/trips/" + trip_id
    }).done(function(cards_unparsed){
      window.location.reload(false);
    });
  })
})


function post_trip () {
  $.ajax({
    method: "POST",
    url: "/lightupon/admin/trips",
    dataType: "json",
    processData: false,
    contentType: "application/json; charset=utf-8",
    data:JSON.stringify({
      "Title":$("#input-trip_title").val(),
      "Description":$("#input-trip_description").val(),
      "ImageURL":$("#input-trip_image_url").val()
    })
  }).done(function(stuff){
    console.log(stuff)
  });
}

})();
</script>
`

const trip_detail_template = `
<style>
  .bold {
    font-weight: bold;
  }

  .block_container {
    border: 1px solid black;
    padding-left:3px;
  }

  .scenes_for_trip {
    width: 400px;
    float: left;
  }

  .add_scene {
    margin-left: 420px;
    width:300px;
  }

  .popular_scenes {
    width:400px;
    position: fixed;
    top: 10;
    left: 740;
  }

  .button {
    border: 1px solid black;
    background-color:#EDA1E7;
    width:50px;
  }

  .cards_details {
    padding-left:20px;
  }

</style>
<html>
  <body>
    <span style="visibility: hidden;" id="tripID">{{.ID}}</span>
    <div style="width: 100%; overflow: hidden;">
      <div class="scenes_for_trip block_container">
              <p class="bold"> {{.Title}} (TripID = {{.ID}})<p>
        <p><img src="{{.ImageUrl}}" height="50" width="150"/></p>
        <p class="bold"> SCENES </p>
        <p class="bold"> SceneOrder - Scene.Name - Coordinates</p>
        {{range $index, $element := .Scenes}}
          <p class="scene-link" id="scene_{{$element.ID}}">
            <a href="/lightupon/admin/scenes/{{.ID}}">
              {{$element.SceneOrder}} - {{$element.Name}} - ({{$element.Latitude}},{{$element.Longitude}})
            </a>
            <span class="delete_scene_link button" id="scene_{{$element.ID}}">delete scene</span>
          </p>
            <!--<p><img src="{{$element.BackgroundUrl}}" style="width:50px;height:50px"/></p>-->
        {{end}}
      </div>

      <div class="add_scene block_container" >
        <p class="bold"> CREATE NEW SCENE </p>
        <p>Scene Title: <input type="text" id="input-scene_title"/></p>
        <p>Scene Order: <input type="text" id="input-scene_order"/></p>
        <p>Latitude: <input type="text" id="input-scene_latitude"/></p>
        <p>Longitude: <input type="text" id="input-scene_longitude"/></p>
        <p>Background Image URL: <input type="text" id="input-scene_backgroundURL"/></p>
        <p class="submit_scene button">Submit</p>
      </div>
      <div class="popular_scenes block_container" >
        <p class="bold"> PRE-MADE SCENES </p>
        <p> In order to add a pre-made scene to your trip, enter your desired SceneOrder for the scene and hit 'Add Scene' </p>
      </div>
    </div>
  </body>
</html>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>
<script>
(function(){

$('.delete_scene_link').each(function(index, element){
  var delete_scene_link = $(element);
  delete_scene_link.on('click', function(element_1){
    var scene_id = delete_scene_link.attr('id').split('_')[1];
    $.ajax({
      method: "DELETE",
      url:"/lightupon/admin/scenes/" + scene_id
    }).done(function(scenes_unparsed){
      window.location.reload(false);
    });
  })
})


$.ajax({
    method: "GET",
    url: "/lightupon/admin/popular_scenes",
    processData: false,
    dataType: "json"
  }).done(function(stuff){
    print_popular_scenes(stuff);
  });

function print_popular_scenes (stuff) {
  var popular_scenes = $('.popular_scenes');
  stuff.forEach(function(element, index){
    // First create the html and append it to the DOM
    var scene = $('<div><p><span class="bold">' + element.Name + '</span>    <input type="text" id="add_popular_scene_' + element.ID + '_scene_order" style="width:30px"/><span class="button" id="add_popular_scene_' + element.ID + '">Add Scene</span></p></div>');
    popular_scenes.append(scene);

    // Then add that sweet sweet onclick action to allow the submitting of the scene
    $('#add_popular_scene_' + element.ID).on("click", function () {
      var scene_order = parseInt($('#add_popular_scene_' + element.ID + '_scene_order').val());
      post_popular_scene(element.ID, scene_order)
      window.location.reload(false);
    })
  })
}

function post_popular_scene (sceneID, sceneOrder) {
  tripID = $("#tripID").html();
  $.ajax({
    method: "POST",
    url: "/lightupon/admin/trips/" + tripID + "/scenes",
    dataType: "json",
    processData: false,
    contentType: "application/json; charset=utf-8",
    data:JSON.stringify({
      "SceneOrder":sceneOrder,
      "ID":sceneID
    })
  }).done(function(stuff){
    console.log(stuff)
  });
}


$('.submit_scene').on('click', function(){
  post_scene();
  window.location.reload(false);
})

// sharknavion
// {"SceneOrder":3, "Name":"new scene", "Latitude":76.567,"Longitude":87.345}
function post_scene () {
  tripID = $("#tripID").html();
  $.ajax({
    method: "POST",
    url: "/lightupon/admin/trips/" + tripID + "/scenes",
    dataType: "json",
    processData: false,
    contentType: "application/json; charset=utf-8",
    data:JSON.stringify({
      "Name":$("#input-scene_title").val(),
      "SceneOrder":parseInt($("#input-scene_order").val()),
      "Latitude":parseFloat($("#input-scene_latitude").val()),
      "Longitude":parseFloat($("#input-scene_longitude").val()),
      "BackgroundUrl":$("#input-scene_backgroundURL").val()
    })
  }).done(function(stuff){
    console.log(stuff)
  });
}
})();
</script>
`

const scene_detail_template = `
<style>
  .bold {
    font-weight: bold;
  }

  .block_container {
    border: 1px solid black;
    padding-left:3px;
  }

  .scenes_for_trip {
    width: 500px;
    border: 1px solid black;
    padding-left:5px;
  }

  .add_card {
    margin-left: 600px;
    margin-top: 10px;
    width:300px;
  }

  .button {
    border: 1px solid black;
    background-color:#EDA1E7;
    width:50px;
  }

</style>
<html>
  <body>
    <span style="visibility: hidden;" id="sceneID">{{.ID}}</span>
    <div style="width: 100%; overflow: hidden;">
      <div class="scenes_for_trip block_container">
        <p class="bold"> SCENE {{.SceneOrder}}<p>
        <p class="bold"> CARDS </p>
        <p class="bold"> CardOrder - Card.Text</p>
        {{range $index, $element := .Cards}}
          <p class="card-link" id="card_{{$element.ID}}">
            {{$element.CardOrder}} - {{$element.Text}}
            <img src="{{$element.ImageURL}}" style="width:100px;height:100px"/>
            <span class="delete_card_link button" id="card_{{$element.ID}}">delete card</span>
          </p>
        {{end}}
      </div>
      <div class="add_card block_container" >
        <p class="bold"> ADD CARD </p>
        <p>Text: <input type="text" id="input-card_text"/></p>
        <p>ImageURL: <input type="text" id="input-card_image_url"/></p>
        <p>CardOrder: <input type="text" id="input-card_card_order"/></p>
        <p>NibID:
          <select id="input-card_nib_id">
            <option value="TextHero">TextHero</option>
            <option value="PictureHero">PictureHero</option>
            <option value="MapHero">MapHero</option>
          </select>
        </p>
        <p class="submit_card button">Submit</p>
      </div>
    </div>
  </body>
</html>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>
<script>
(function(){

$('.submit_card').on('click', function(){
  post_card();
  window.location.reload(false);
})

$('.delete_card_link').each(function(index, element){
  var delete_card_link = $(element);
  delete_card_link.on('click', function(element_1){
    var card_id = delete_card_link.attr('id').split('_')[1];
    console.log(card_id)
      
    $.ajax({
      method: "DELETE",
      url:"/lightupon/admin/cards/" + card_id
    }).done(function(cards_unparsed){
      window.location.reload(false);
    });
  })
})


function post_card () {
  sceneID = parseInt($('#sceneID').html());
  $.ajax({
    method: "POST",
    url: "/lightupon/admin/scenes/" + sceneID + "/cards",
    dataType: "json",
    processData: false,
    contentType: "application/json; charset=utf-8",
    data:JSON.stringify({
      "Text":$("#input-card_text").val(),
      "ImageURL":$("#input-card_image_url").val(),
      "SceneID":sceneID,
      "CardOrder":parseInt($("#input-card_card_order").val()),
      "NibID":$("#input-card_nib_id").val()
    })
  }).done(function(stuff){
    console.log(stuff)
  });
}

})();
</script>
`