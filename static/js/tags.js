(function() {
	'use strict';

	// Copy|Delete|manage tags inputs
	document.addEventListener('DOMContentLoaded', function() {

		/*
		** Template
		*
		*	<input class="tags-add-input" /> -- Main input
		*	<div class="tags-container"> -- Tags container
		** 	<button class="tags-add-button"> -- If you also want an 'add' button (Cf Khaled)
		*/

		var tagsContainer = document.getElementsByClassName('form-tag-list');

		if (tagsContainer.length === 0)
			return ;

		for (var i = 0; i < tagsContainer.length; i++) {

			new FormTag(tagsContainer[i]);
		}

		function FormTag(element) {

			this.tagsContainer = element;

			var tags = this.tagsContainer.getElementsByClassName('form-tag-item')
			for (var i = 0; i < tags.length; i++) {
				tags[i].getElementsByTagName('button')[0].onclick = removeTag;
			}

			this.tagsInput = this.tagsContainer.getElementsByClassName('form-tag-add')[0];

			var self = this;
			this.tagsInput.onkeypress = function(event) {
				if (event.which === 13 || event.keyCode === 13) {
					event.preventDefault();
					return self.addTags(event);
				}
			};
		}

		FormTag.prototype.addTags = function(event) {

			var value = this.tagsInput.value;
			if (value === undefined || value.length === 0)
				return;
			this.tagsContainer.insertBefore(this.buildElem(value), this.tagsInput);
			this.tagsInput.value = '';
		}

		FormTag.prototype.buildElem = function(value) {

			var span = document.createElement('span');
			var input = document.createElement('input');
			var removeButton = document.createElement('button');

			span.classList.add('form-tag-item');

			input.setAttribute('value', value);
			input.setAttribute('name', this.tagsInput.getAttribute('data-name')+'[]');
			input.setAttribute('type', 'text');
			input.setAttribute('readonly', '');

			removeButton.innerHTML = 'x';
			removeButton.onclick = removeTag;

			span.appendChild(input);
			span.appendChild(removeButton);
			return span;
		}

		function removeTag(event) {
			event.preventDefault();
			var p = this.parentNode;
			p.parentNode.removeChild(p);
		}

	});

})();
