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
        <p class="bold"> TRIP {{.ID}}<p>
        <p class="bold"> SCENES </p>
        <p class="bold"> SceneOrder / SceneID / Scene.Name</p>
        {{range $index, $element := .Scenes}}
          <p class="scene-link" id="scene_{{$element.ID}}">
            <a href="/lightupon/admin/scenes/{{.ID}}">
              {{$element.SceneOrder}} / {{$element.ID}} / {{$element.Name}}
            </a>
          </p>
        {{end}}
      </div>

      <div class="add_scene block_container" >
        <p class="bold"> ADD SCENE </p>
        <p>Scene Title: <input type="text" id="input-scene_title"/></p>
        <p>Scene Order: <input type="text" id="input-scene_order"/></p>
        <p>Latitude: <input type="text" id="input-scene_latitude"/></p>
        <p>Longitude: <input type="text" id="input-scene_longitude"/></p>
        <p class="submit_scene button">Submit</p>
      </div>
    </div>
  </body>
</html>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>
<script>
(function(){

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
      "Longitude":parseFloat($("#input-scene_longitude").val())
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
        <p class="bold"> CardOrder / CardID / Card.Text</p>
        {{range $index, $element := .Cards}}
          <p class="card-link" id="card_{{$element.ID}}">
            {{$element.CardOrder}} / {{$element.ID}} / {{$element.Text}}
            <span class="delete_card_link button" id="card_{{$element.ID}}">delete card</span>
          </p>
        {{end}}
      </div>
      <div class="add_card block_container" >
        <p class="bold"> ADD CARD </p>
        <p>Text: <input type="text" id="input-card_text"/></p>
        <p>ImageURL: <input type="text" id="input-card_image_url"/></p>
        <p>SceneID: <input type="text" id="input-card_scene_id"/></p>
        <p>CardOrder: <input type="text" id="input-card_card_order"/></p>
        <p>NibID:
          <select id="input-card_nib_id">
            <option value="textHero">textHero</option>
            <option value="pictureHero">pictureHero</option>
            <option value="mapHero">mapHero</option>
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
  sceneID = parseInt($("#input-card_scene_id").val());
  $.ajax({
    method: "POST",
    url: "/lightupon/admin/scenes/" + sceneID + "/cards",
    dataType: "json",
    processData: false,
    contentType: "application/json; charset=utf-8",
    data:JSON.stringify({
      "Text":$("#input-card_text").val(),
      "ImageURL":$("#input-card_image_url").val(),
      "SceneID":parseInt($("#input-card_scene_id").val()),
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