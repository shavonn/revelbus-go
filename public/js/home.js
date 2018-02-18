$(document).ready(function(){
    $(".owl-carousel").owlCarousel({
        items: 1,
        autoplay: true,
        loop: true,
        nav:true,
        dots: true
    });
    
    var $gallery = $('.gallery a').simpleLightbox();
});