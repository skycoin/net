import { Component, OnInit, ViewEncapsulation, Input } from '@angular/core';

@Component({
  selector: 'app-im-history-message',
  templateUrl: './im-history-message.component.html',
  styleUrls: ['./im-history-message.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ImHistoryMessageComponent implements OnInit {
  @Input() type = 'other';
  @Input() msg = '';
  @Input() from = '';
  constructor() { }

  ngOnInit() {
  }

}
