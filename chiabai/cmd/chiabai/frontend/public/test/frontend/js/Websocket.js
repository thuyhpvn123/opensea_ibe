document.addEventListener('DOMContentLoaded', async() => {
  try{
    console.log("begin")
   initApp()
  //  loadImage()
  }catch(e){
   console.log(e.message)
  }
});
var files =[]
var fromAddress,prikey
var flagUpload=1
var $resultCall = document.getElementById('result-call');

var dataCreateCall = JSON.stringify(
  {
    'from':   "45c75cfb8e20a8631c134555fa5d61fcf3e602f2",
    'priKey': "36e1aa979f98c7154fb2491491ec044ccac099651209ccfbe2561746dbe29ebb",
    'to':    "227e5cba2e6f953512a58649cd661643442bf096",
    amount:  "",
    "function-name":"domain",
    gas:1000000,
    gasPrice:10,
    timeUse:1000,
    relatedAddresses:[],
  }
);
var tokenURI=JSON.stringify(
  {
    "internalType": "string",
    "name": "tokenURI",
    "type": "string",
    "value": "http://localhost:2000/getNft/",
  },
)
var price=JSON.stringify(
  {
    "internalType": "uint256",
    "name": "price",
    "type": "uint256",
    "value": 1000000000000000000,
  }
)

var Inputs=[tokenURI,price]
var relatedAdd1 ={
  "address":"45c75cfb8e20a8631c134555fa5d61fcf3e602f2"
}
var dataCall = JSON.stringify(
  {
    'from':   "45c75cfb8e20a8631c134555fa5d61fcf3e602f2",
    'priKey': "36e1aa979f98c7154fb2491491ec044ccac099651209ccfbe2561746dbe29ebb",
    'to':    "227e5cba2e6f953512a58649cd661643442bf096",
    amount:            "2540BE400",
    "function-name":"createToken",
    inputArray:Inputs,
    gas:1000000,
    gasPrice:10,
    timeUse:1000,
    relatedAddresses:[relatedAdd1],
  }
);


//http://localhost:3000/api/v1/test/template/
const initApp = async()=>{
//upload
  var $upload = document.getElementById('file');
  var $createResult = document.getElementById('create-result');
  var $createResultCall = document.getElementById('create-result-call');

  var $connectNode = document.getElementById('register');
  var $sendJSON = document.getElementById('information');

  $upload.addEventListener("change", async(e) => { 
    e.preventDefault();
    files = e.target.files;
   await  getImage(files[0]);
  });
 

  $sendJSON.addEventListener('submit', async(e) => {

    e.preventDefault()
    console.log("create nft")
    var name, description,traitType,value, maxValue,flag =1,
    name = $('#name').val()
    description = $('#description').val()
    traitType = $('#trait_type').val()
    value = $('#value').val()
    maxValue = $('#max_value').val()

    if( name ==''){
      flag=0
      $('.error_name').html("Please type name of the nft")
    }else{
      $('.error_name').html("")
    }
    if( description ==''){
      flag=0
      $('.error_description').html("Please type description of the nft")
    }else{
      $('.error_description').html("")
    }
    if( value ==''){
      flag=0
      $('.error_value').html("Please type value of attributes")
    }else{
      $('.error_value').html("")
    }
    if( traitType ==''){
      flag=0
      $('.error_trait_type').html("Please type trait_type of attributes")
    }else{
      $('.error_trait_type').html("")
    }
    if( maxValue ==''){
      flag=0
      $('.error_max_value').html("Please type max_value of attributes")
    }else{
      $('.error_max_value').html("")
    }

    //create new nft
    if(flag==1 && flagUpload==0 ){
      try{
        console.log("begin to create nft")
            // Creating a XHR object
    let xhr = new XMLHttpRequest();
    let url = "/createNft";

    // open a connection
    xhr.open("POST", url, true);

    // Set the request header i.e. which type of content you are sending
    xhr.setRequestHeader("Content-Type", "application/json");

    // Create a state change callback
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {

            // Print received data from server
            var data=xhr.responseText;
            var jsonResponse = JSON.parse(data);
            var tokenid = jsonResponse["data"]["tokenid"]
            $createResultCall.innerHTML =` json in database with TokenID: ${tokenid}` ;
        }
    };
    var atr1={
      "trait_type": traitType,
      "value": parseInt(value),
      "max_value": parseInt(maxValue)
    }

    var attributeArr =[atr1]
    // Converting JSON data to string
    var data = JSON.stringify(
      { 
        "name": name, 
        "description": description,
        "attributes":attributeArr
      }
      );
    console.log("data là:",data)

    xhr.send(data);

        // await createNftInfo()
        // getImage(files[0]);
      createCall();
      }catch{
        console.log(e)
        $createResultCall.innerHTML = `Ooops... there was an error while trying to create a new nft`;
      }
    }else{
      $createResultCall.innerHTML = `Please choose file or enter needed input`;
    }

  })

  $connectNode.addEventListener('submit', async(e) => {

    e.preventDefault()
    console.log("connect node")
    var flag =1,
    fromAddress = $('#address').val()
    prikey = $('#prikey').val()

    if( fromAddress ==''){
      flag=0
      $('.error_address').html("Please type address of the account")
    }else{
      $('.error_address').html("")
    }
    if( prikey ==''){
      flag=0
      $('.error_prikey').html("Please type private key of the account")
    }else{
      $('.error_prikey').html("")
    }

    if(flag==1  ){
      try{
        console.log("begin to connect node")
        // Creating a XHR object
        let xhr = new XMLHttpRequest();
        let urlConnect="/connectWallet"

        // open a connection
        xhr.open("POST", urlConnect, true);

        // Set the request header i.e. which type of content you are sending
        xhr.setRequestHeader("Content-Type", "application/json");

        // Create a state change callback
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4 && xhr.status === 200) {

                // Print received data from server
                var data=xhr.responseText;
                var jsonResponse = JSON.parse(data);
                var addr = jsonResponse["data"]["address"]
                $createResult.innerHTML =` connect wallet ${addr} successfully` ;
                $('#wallet-id-user').html(`${addr}`) ;

                // $createResult.innerHTML =data ;

            }
        };
        // Converting JSON data to string
        var data = JSON.stringify(
          { 
            "address": fromAddress, 
            "privatekey": prikey,
          }
          );
        console.log("data là:",data)

        xhr.send(data);
      }catch{
        console.log(e)
        $createResult.innerHTML = `Ooops... there was an error while trying to register wallet`;
      }
    }
  })
}

const createCall = () => {
  // Creating a XHR object
  let xhr = new XMLHttpRequest();

  // open a connection
  xhr.open("POST", "/call", true);

  // Set the request header i.e. which type of content you are sending
  xhr.setRequestHeader("Content-Type", "application/json");

  // Create a state change callback
  xhr.onreadystatechange = function () {
      if (xhr.readyState === 4 && xhr.status === 200) {

          // Print received data from server
          var data=xhr.responseText;
          console.log("ket qua call là:",data)
          // var jsonResponse = JSON.parse(data);
          // var result = jsonResponse["result"]

          $resultCall.innerHTML =`you just created a new nft and listed for sale ` ;

      }
  };
  // Converting JSON data to string

  console.log("data là:",dataCall)

  xhr.send(dataCall);
};


const getImage = async(file) => {
  // Creating a XHR object
  let xhr = new XMLHttpRequest();
  var fd = new FormData();
  fd.append("file", file);
  await xhr.open("POST", '/upload',true);
   xhr.onreadystatechange = function () {
     if (xhr.readyState == 4&& xhr.status === 200) {
       console.log(xhr.responseText);
        showImage();
      }
    };
    xhr.send(fd);
};
const ShowProduct=(a)=>{
  var Pshow = ` <div class="col-xs-6 col-sm-6 col-md-6 col-lg-6" style="margin-top:50px">
      <div id="product-1">
        <a href="#" class="thumbnail" width="500px" height="500px">
          <div class="image-1">
          <img src="../../uploads/${a}" width="200px" height="200px"></img>
          </div>
        </a>             
      </div>
    </div>`  
  $('.product').append(Pshow)     
}

const showImage = async ()=>{
  if (flagUpload==1){
    ShowProduct(files[0].name);
    flagUpload =0;
  }else{
    const element = document.getElementById('product-1');
    element.remove()
    ShowProduct(files[0].name);
    flagUpload =0
  }
  console.log("def");
}

  // var loadImage = ()=>{
  //   fetch('/getNft')
  //   .then(function(response) {
  //     if (!response.ok) {
  //       throw Error(response.statusText);
  //     }
  //     // Read the response as json.
  //     return response.json();
  //   })
  //   .then(function(responseAsJson) {
  //     // Do stuff with the JSON
  //       console.log("responseAsJson:",responseAsJson)
      
  //   })
  //   .catch(function(error) {
  //     console.log('Looks like there was a problem: \n', error);
  //   });
      
  // }

 // var loadImage = (id)=>{
//   // Creating a XHR object
//   let xhr = new XMLHttpRequest();
//   var param = "tokenid="+id
//   // open a connection
//   xhr.open("GET", "/getNft"+"?"+param, true);

//   // Set the request header i.e. which type of content you are sending
//   xhr.setRequestHeader("Content-Type", "application/json");

//   // Create a state change callback
//   xhr.onreadystatechange = function () {
//       if (xhr.readyState === 4 && xhr.status === 200) {

//           // Print received data from server
//           var data=xhr.responseText;
//           console.log("ket qua call là:",data)
//           var jsonResponse = JSON.parse(data);
//           var result = jsonResponse["result"]

//           $createResultCall.innerHTML =`result: ${result}` ;
//           console.log(" result là:",result)

//       }
//     };

//     xhr.send(null);

// }