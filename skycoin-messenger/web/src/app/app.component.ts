import { Component, OnInit, ViewChild, OnDestroy } from '@angular/core';
import { SocketService } from '../providers';
import { ImRecentItemComponent } from '../components';

@Component({
  selector: 'app-im',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnDestroy {
  constructor(public ws: SocketService) {
  }

  ngOnDestroy() {
    this.ws.socket.unsubscribe();
  }
}
