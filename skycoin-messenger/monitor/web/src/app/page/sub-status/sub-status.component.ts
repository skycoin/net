import { Component, OnInit, ViewEncapsulation, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { environment as env } from '../../../environments/environment';
import { ApiService, NodeServices, App, Transports } from '../../service';
import { MatSnackBar, MatDialog } from '@angular/material';
import { DataSource } from '@angular/cdk/collections';
import { Observable } from 'rxjs/Observable';
import { UpdateCardComponent } from '../../components/update-card/update-card.component';
import 'rxjs/add/observable/of';

@Component({
  selector: 'app-sub-status',
  templateUrl: './sub-status.component.html',
  styleUrls: ['./sub-status.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class SubStatusComponent implements OnInit, OnDestroy {
  sshColumns = ['index', 'key', 'del'];
  displayedColumns = ['index', 'key', 'app'];
  transportColumns = ['index', 'fromNode', 'fromApp', 'toNode', 'toApp'];
  appSource: SubStatusDataSource = null;
  sshSource: SubStatusDataSource = null;
  transportSource: SubStatusDataSource = null;
  key = '';
  power = '';
  transports: Array<Transports> = [];
  status: NodeServices = null;
  apps: Array<App> = [];
  task = null;
  isManager = false;
  socketColor = 'close-status';
  sshColor = 'close-status';
  dialogTitle = '';
  sshTextarea = '';
  sshAllowNodes = [];
  socksTextarea = '';
  sshConnectKey = '';
  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private api: ApiService,
    private snackBar: MatSnackBar,
    private dialog: MatDialog) {
  }

  ngOnInit() {
    this.route.queryParams.subscribe(params => {
      this.key = params.key;
      this.startTask();
      this.power = 'warn';
      this.isManager = env.isManager;
    });
  }
  ngOnDestroy() {
    this.close();
  }
  delAllowNode(ev: Event, index: number) {
    this.sshAllowNodes.splice(index, 1);
    this.api.runSSHServer(this.status.addr, this.sshAllowNodes.toString()).subscribe(result => {
      if (result) {
        console.log('set ssh result:', result);
        this.sshTextarea = '';
        this.init();
      }
    });
  }
  setSSH(ev: Event) {
    let data = '';
    if (this.sshTextarea.trim()) {
      data = this.sshAllowNodes + ',' + this.sshTextarea.trim();
    } else {
      data = this.sshAllowNodes.toString();
    }
    this.api.runSSHServer(this.status.addr, data).subscribe(result => {
      if (result) {
        console.log('set ssh result:', result);
        this.sshTextarea = '';
        this.init();
      }
    });
  }
  checkUpdate(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.dialog.open(UpdateCardComponent, {
      panelClass: 'update-panel',
      width: '370px',
      disableClose: true
    });
  }
  refresh(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.init();
    this.snackBar.open('Refreshed', 'Dismiss', {
      duration: 3000,
      verticalPosition: 'top',
      extraClasses: ['bg-success']
    });
  }
  runSocketServer(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.api.runSockServer(this.status.addr).subscribe(isOk => {
      if (isOk) {
        this.init();
        console.log('start socket server');
      }
    });
  }
  runSSHServer(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    this.api.runSSHServer(this.status.addr).subscribe(isOk => {
      if (isOk) {
        this.init();
        console.log('start ssh server');
      }
    });
  }
  reboot(ev: Event) {
    ev.stopImmediatePropagation();
    ev.stopPropagation();
    ev.preventDefault();
    console.log('reboot');
    this.api.reboot(this.status.addr).subscribe(isOk => {
      if (isOk) {
        if (this.task) {
          this.close();
        }
        this.startTask();
      }
    });
  }
  inputKeys(ev: Event, content: any) {
    this.dialog.open(content, {
      width: '450px'
    });
  }
  openSettings(ev: Event, content: any, title: string) {
    if (!title) {
      return;
    }
    this.sshSource = new SubStatusDataSource(this.findService('ssh').allow_nodes);
    this.dialogTitle = title;
    this.dialog.open(content, {
      width: '800px',
    });
  }
  findService(search: string): App {
    const result = this.status.apps.find(el => {
      const tmp = el.attributes.find(attr => {
        return search === attr;
      });
      return search === tmp;
    });
    return result;
  }
  startTask() {
    this.init();
    this.task = setInterval(() => {
      this.init();
    }, 15000);
  }
  close() {
    clearInterval(this.task);
    this.task = null;
  }
  isExist(search: string) {
    const result = this.status.apps.find(el => {
      const tmp = el.attributes.find(attr => {
        return search === attr;
      });
      return search === tmp;
    });
    return result !== undefined && result !== null;
  }
  setServiceStatus() {
    this.sshColor = 'close-status';
    this.socketColor = 'close-status';
    if (this.status.apps) {
      this.sshAllowNodes = this.findService('ssh').allow_nodes;
      this.sshSource = new SubStatusDataSource(this.sshAllowNodes);
      if (this.isExist('ssh')) {
        this.sshColor = 'mat-primary';
      }
      if (this.isExist('socks')) {
        this.socketColor = 'mat-primary';
      }
    }
  }
  fillTransport() {
    if (env.isManager && this.status.addr) {
      this.api.getTransport(this.status.addr).subscribe((allTransports: Array<Transports>) => {
        if (allTransports && allTransports.length > 0) {
          this.transports = allTransports;
          this.transportSource = new SubStatusDataSource(allTransports);
        }
      });
    }
  }
  fillApps() {
    if (this.status.apps) {
      this.appSource = new SubStatusDataSource(this.status.apps);
      this.setServiceStatus();
    } else {
      if (env.isManager) {
        this.api.getApps(this.status.addr).subscribe((apps: Array<App>) => {
          this.status.apps = apps;
          this.appSource = new SubStatusDataSource(this.status.apps);
          this.setServiceStatus();
        });
      }
    }
  }
  init() {
    const data = new FormData();
    data.append('key', this.key);
    this.api.getNodeStatus(data).subscribe((nodeServices: NodeServices) => {
      if (nodeServices) {
        this.status = nodeServices;
        this.fillTransport();
        this.fillApps();
      }
    });
  }
}
export class SubStatusDataSource extends DataSource<any> {
  size = 0;
  constructor(private apps: any) {
    super();
  }
  connect(): Observable<any> {
    return Observable.of(this.apps);
  }

  disconnect() { }
}
