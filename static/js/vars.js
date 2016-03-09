(function() {
  'use strict';

  document.addEventListener('DOMContentLoaded', function () {
    //Form "vars" input
    var formVar = document.getElementsByClassName('form-var');
    for (var i = 0; i < formVar.length; i++) {
        new FormVar(formVar[i]);
    }

    function FormVar(element) {

        var self = this;

        self.form = element;

        // Form template
        self.template = self.form.getElementsByClassName('form-var-template')[0];

        // Form var-removefire
        var varRemoves = self.form.getElementsByClassName('form-var-remove');
        for (var i = 0; i < varRemoves.length; i++) {
            varRemoves[i].addEventListener('click', function() {
                this.parentNode.parentNode.removeChild(this.parentNode);
            });
        }

        // Add button
        var actionButton = self.form.getElementsByClassName('form-var-add')[0];
        if (actionButton !== null)
          //Append new group of inputs when "+" button is clicked
          actionButton.addEventListener('click', function(e) {
            e.preventDefault();

            // Clone template
            var node = self.template.cloneNode(true);
            node.classList.remove('is-hidden');
            node.classList.remove('form-var-template');
            node.classList.add('form-var-item');

            var inputs = node.getElementsByTagName('input');
            for (var i = 0; i < inputs.length; i++) {
                inputs[i].setAttribute('name', inputs[i].getAttribute('data-name'));
            }
            var selects = node.getElementsByTagName('select');
            for (var i = 0; i < selects.length; i++) {
                selects[i].setAttribute('name', selects[i].getAttribute('data-name'));
            }

            // Setup remove button
            var btn = node.getElementsByClassName('form-var-remove')[0];
            if (btn !== null ) {
                btn.addEventListener('click', function(e) {
                    e.preventDefault();
                    this.parentNode.parentNode.removeChild(this.parentNode);
                });
            }

            // Add node
            self.form.insertBefore(node, actionButton);
        });
    }
  });
})();
