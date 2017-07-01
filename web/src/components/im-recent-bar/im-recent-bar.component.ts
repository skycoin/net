import { Component, OnInit, ViewEncapsulation, Input, ViewChildren, QueryList } from '@angular/core';
import { ImRecentItemComponent } from '../im-recent-item/im-recent-item.component';

@Component({
  selector: 'app-im-recent-bar',
  templateUrl: './im-recent-bar.component.html',
  styleUrls: ['./im-recent-bar.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class ImRecentBarComponent implements OnInit {
  @ViewChildren(ImRecentItemComponent) items: QueryList<ImRecentItemComponent>;
  @Input() list = [];
  constructor() { }

  ngOnInit() {
  }

  selectItem(item: ImRecentItemComponent) {
    const tmp = this.items.filter((el) => {
      return el.name !== item.name;
    });
    tmp.forEach(el => {
      el.active = false;
    });
  }
}
