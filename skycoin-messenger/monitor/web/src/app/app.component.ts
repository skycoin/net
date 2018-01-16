import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  listitems = [{ text: 'Node List', link: '/', icon: 'list' }, { text: 'Update Password', link: '/updatePass', icon: 'lock' }];
  option = { selected: true };
  constructor() { }
}
