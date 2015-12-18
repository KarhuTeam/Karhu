(function() {
  'use strict';

  function VarInput() {
    var toReturn = document.createElement('div');
    toReturn.classList.add('var-unit-container');
    toReturn.innerHTML = "<input name=\"varKeys[]\" type=\"text\" placeholder=\"Key\" /><input name=\"varValues[]\" type=\"text\" placeholder=\"Value\" />";

    return toReturn;
  }

  document.addEventListener('DOMContentLoaded', function () {
    //Form "vars" input
    var formVar = document.getElementsByClassName('form-var')[0];
    if (formVar !== null) {
        var actionButton = document.getElementsByClassName('add-var-button')[0];
        if (actionButton !== null)
          //Append new group of inputs when "+" button is clicked
          actionButton.addEventListener('click', function(e) {
            e.preventDefault();
            formVar.appendChild(new VarInput());
        });
    }
  });
})();
