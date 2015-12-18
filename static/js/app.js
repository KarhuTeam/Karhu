(function() {
  'use strict';

  document.addEventListener('DOMContentLoaded', function () {
    //Navbar toggle display
    var navbarBtn = document.getElementsByClassName('navbar-btn')[0];
    if (navbarBtn !== null)
      navbarBtn.addEventListener('click', function(e) {
        e.preventDefault();
        document.getElementsByClassName('header')[0].classList.toggle('is-open');
      });
  });
})();
