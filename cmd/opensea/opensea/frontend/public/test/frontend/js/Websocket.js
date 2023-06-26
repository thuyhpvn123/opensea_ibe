document.addEventListener('DOMContentLoaded', async() => {
  try{
   await initWeb3()
 //   await getContractInfo()
   await initContract();
   await Display()
   await loadData()
   await initApp(); 
   await Edit()
   await Start()
 //   await initContractInfura()
  }catch(e){
   console.log(e.message)
  }
});
const $createResult = document.getElementById('create-result');
const Create= async()=>{
  $(document).on('submit','#information',async(e)=>{
    e.preventDefault()
    var name, des,fee,btnCre, flag =1,
    name = $('#name').val()
    des = $('#description').val()
    fee = $('#fee').val()
    btnCre = $('#information').val()

    if( name ==''){
      flag=0
      $('.error_name').html("Please type name of the nft")
    }else{
      $('.error_name').html("")
    }
    if( des ==''){
      flag=0
      $('.error_des').html("Please type description of the nft")
    }else{
      $('.error_des').html("")
    }
    if( des ==''){
      flag=0
      $('.error_fee').html("Please type fee to mint nft")
    }else{
      $('.error_fee').html("")
    }
	
    //create new nft
    if(flag==1  ){
      try{
        var givenPrice,idBtn0;
        givenPrice = parseInt($(`#givenPrice-${a}`).val())
      }catch{
        console.log(e)
        $createResult.innerHTML = `Ooops... there was an error while trying to create a new nft`;
      }
    }
  })
}