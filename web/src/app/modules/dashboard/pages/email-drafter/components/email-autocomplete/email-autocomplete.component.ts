import { NgIf } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-email-autocomplete',
  imports: [NgIf],
  templateUrl: './email-autocomplete.component.html',
  styleUrl: './email-autocomplete.component.css'
})
export class EmailAutocompleteComponent {
  @Input() suggestions: string[] = [];
  @Output() suggestionSelected = new EventEmitter<string>();

  selectSuggestion(suggestion: string) {
    this.suggestionSelected.emit(suggestion);
  }
}
