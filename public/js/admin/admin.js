$(document).ready(function(){
    ClassicEditor
        .create(document.querySelector('.editor'));

    flatpickr('.datetime_field', {
        enableTime: true,
        dateFormat: "Y-m-d H:i",
        altInput: true,
        altFormat: "F j, Y h:iK",
    });
});