import { Component, OnInit, ViewEncapsulation } from '@angular/core';

@Component({
  selector: 'app-im-history-view',
  templateUrl: './im-history-view.component.html',
  styleUrls: ['./im-history-view.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ImHistoryViewComponent implements OnInit {
  testList = [1, 2, 3, 4, 5, 6, 7];
  constructor() { }
  ngOnInit() {
  }

}
