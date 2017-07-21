import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { FormControl, Validators } from '@angular/forms';

@Component({
  // tslint:disable-next-line:component-selector
  selector: 'create-chat-dialog',
  templateUrl: './create-chat-dialog.component.html',
  styleUrls: ['./create-chat-dialog.component.scss'],
  encapsulation: ViewEncapsulation.None
})

export class CreateChatDialogComponent implements OnInit {
  keys = '';
  keyFormControl = new FormControl('', [
    Validators.required]);
  constructor() { }

  ngOnInit() {
  }
}
