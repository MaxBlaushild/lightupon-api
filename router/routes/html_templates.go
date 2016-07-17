package routes

const trips_list_template = `
<style>
  .bold {
    font-weight: bold;
  }

  .scenes_for_trip {
    width: 500px;
    border: 1px solid black;
    padding-left:5px;

  }

</style>
<html>
  <body>
    <div class="all-trips">
      <p class="bold"> YOUR TRIPS<p>
        {{range $index, $element := .}}
          <p class="bold">{{$element.Title}}</p>
          <p>{{$element.Description}}</p>
          <p><img src="{{$element.ImageUrl}} height="200" width="200"/></p>
        {{end}}
    </div><br><br>
  </body>
</html>
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

  .add_card {
    margin-left: 840px;
    width:300px;
  }

  .submit_button {
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
            {{$element.SceneOrder}} / {{$element.ID}} / {{$element.Name}}
          </p>
        {{end}}
      </div>
      <div class="add_scene block_container" >
        <p class="bold"> ADD SCENE </p>
        <p>Scene Title: <input type="text" id="input-scene_title"/></p>
        <p>Scene Order: <input type="text" id="input-scene_order"/></p>
        <p>Latitude: <input type="text" id="input-scene_latitude"/></p>
        <p>Longitude: <input type="text" id="input-scene_longitude"/></p>
        <p class="submit_scene submit_button">Submit</p>
      </div>

      <div class="add_card block_container" >
        <p class="bold"> ADD CARD </p>
        <p>Text: <input type="text" id="input-card_text"/></p>
        <p>ImageURL: <input type="text" id="input-card_image_url"/></p>
        <p>SceneID: <input type="text" id="input-card_scene_id"/></p>
        <p>CardOrder: <input type="text" id="input-card_card_order"/></p>
        <p>NibID: <input type="text" id="input-card_nib_id"/></p>
        <p class="submit_card submit_button">Submit</p>
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

$('.submit_card').on('click', function(){
  post_card();
  window.location.reload(false);
})

$('.scene-link').each(function(index, element){
  var scene_link = $(element);
  scene_link.on('click', function(element_1){
    var scene_id = scene_link.attr('id').split('_')[1];
    $.ajax({
      method: "GET",
      // url:"http://localhost:5000/lightupon/admin/scenes/" + scene_id + "/cards",
      url: "http://45.55.160.25/lightupon/admin/scenes/" + scene_id + "/cards",
      datatype:"json"
    }).done(function(cards_unparsed){
      var cards = JSON.parse(cards_unparsed);
      var html_to_append = '<div class="cards_details" id="cards_for_scene_' + scene_id + '"><span class="bold">CARDS</span>';
      for (i=0; i<cards.length; i++) {
        html_to_append += '<p>' + i +  ' / ' + cards[i]["NibID"] + ' / ' +  ' / ' + cards[i]["Text"] + '</p>'
      }
      html_to_append += '</div>';
      scene_link.after(html_to_append);
    });
  })
})

function post_card () {
  sceneID = parseInt($("#input-card_scene_id").val());
  $.ajax({
    method: "POST",
    // url: "http://localhost:5000/lightupon/admin/scenes/" + sceneID + "/cards_post",
    url: "http://45.55.160.25/lightupon/admin/scenes/" + sceneID + "/cards_post",
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

// sharknavion
// {"SceneOrder":3, "Name":"new scene", "Latitude":76.567,"Longitude":87.345}
function post_scene () {
  tripID = $("#tripID").html();
  $.ajax({
    method: "POST",
    // url: "http://localhost:5000/lightupon/admin/trips/" + tripID + "/scenes_post",
    url: "http://45.55.160.25/lightupon/admin/trips/" + tripID + "/scenes_post",
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