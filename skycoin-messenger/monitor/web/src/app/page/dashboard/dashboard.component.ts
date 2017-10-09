import { Component, OnInit, ViewEncapsulation, OnDestroy } from '@angular/core';
import { ApiService, ConnData, ConnsResponse } from '../../service';
import { DataSource } from '@angular/cdk/collections';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';
import { MdSnackBar } from '@angular/material';
@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class DashboardComponent implements OnInit, OnDestroy {
  displayedColumns = ['index', 'type', 'status', 'key', 'send', 'receive', 'seen'];
  dataSource: ExampleDataSource = null;
  dataSize = 0;
  refreshTask = null;
  constructor(private api: ApiService, private snackBar: MdSnackBar) { }
  ngOnInit() {
    this.refresh();
    this.refreshTask = setInterval(() => {
      this.refresh();
    }, 10000);
  }
  ngOnDestroy() {
    this.close();
  }
  TrackByKey(index, item) {
    return item ? item.index : undefined;
  }
  status(ago: number) {
    const now = new Date().getTime() / 1000;
    return (now - ago) < 120;
  }
  refresh(ev?: Event) {
    if (ev) {
      ev.stopImmediatePropagation();
      ev.stopPropagation();
      ev.preventDefault();
    }
    this.dataSource = new ExampleDataSource(this.api);
    if (ev) {
      this.snackBar.open('Refresh Successful', 'Dismiss', {
        duration: 2000,
        verticalPosition: 'top'
      });
    }
  }
  close() {
    clearInterval(this.refreshTask);
  }
}
export class ExampleDataSource extends DataSource<any> {
  size = 0;
  constructor(private api: ApiService) {
    super();
  }
  connect(): Observable<ConnData[]> {
    return this.api.getAllNode().map((res: ConnsResponse) => {
      this.size = res.conns.length;
      const connDatas: Array<ConnData> = [];
      res.conns.forEach((v, i) => {
        connDatas.push({
          index: i + 1,
          key: v.key,
          type: v.type,
          send_bytes: v.send_bytes,
          received_bytes: v.received_bytes,
          last_ack_time: v.last_ack_time,
          start_time: v.start_time
        });
      });
      return connDatas;
    });
  }

  disconnect() { }
}
