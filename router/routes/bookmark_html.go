package routes


const bookmark_template = `
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

  .unshared {
    background-color:yellow;
  }

</style>
<html>
  <body>
    <div style="width: 100%; overflow: hidden;">
      <div class="all-trips block_container">
        <p class="bold"> YOUR BOOKMARKS<p>
        <table>
          {{range $index, $element := .}}
          <tr>
            <td>
              <p class="bold"><a href="{{$element.URL}}">{{$element.Title}}</a></p>
            </td>
            <td>
              <p id="bookmark/{{$element.ID}}/mark_semiprivate" class="share_button unshared">&nbsp&nbsp&nbsp</p>
            </td>
            
          </tr>
          {{end}}
        </table>
      </div>
    </div>
  </body>
</html>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>
<script>
(function(){

// $(".share_button").toggleClass("unshared");

$('.share_button').each(function(index, element){
  var element = $(element);
  var lightuponURI = $(element).attr('id');
  element.on('click', function(burt){
    $.ajax({
      method: "PUT",
      url:"/lightupon/admin/" + lightuponURI
    }).done(function(cards_unparsed){
      // window.location.reload(false);
      alert("done");
    });
  })
})



// $('.submit_trip').on('click', function(){
//   post_trip();
//   window.location.reload(false);
// })

// $('.delete_trip').each(function(index, element){
//   var delete_trip_element = $(element);
//   delete_trip_element.on('click', function(element_1){
//     var trip_id = delete_trip_element.attr('id').split('_')[1];
//     console.log(trip_id)
      
//     $.ajax({
//       method: "DELETE",
//       url:"/lightupon/admin/trips/" + trip_id
//     }).done(function(cards_unparsed){
//       window.location.reload(false);
//     });
//   })
// })


// function post_trip () {
//   $.ajax({
//     method: "POST",
//     url: "/lightupon/admin/trips",
//     dataType: "json",
//     processData: false,
//     contentType: "application/json; charset=utf-8",
//     data:JSON.stringify({
//       "Title":$("#input-trip_title").val(),
//       "Description":$("#input-trip_description").val(),
//       "ImageURL":$("#input-trip_image_url").val()
//     })
//   }).done(function(stuff){
//     console.log(stuff)
//   });
// }

})();
</script>
`