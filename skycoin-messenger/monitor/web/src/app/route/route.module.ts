import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { DashboardComponent, SubStatusComponent, LoginComponent, UpdatePassComponent } from '../page';


const home = {
  path: '',
  component: DashboardComponent,
  pathMatch: 'full'
};
const node = {
  path: 'node',
  component: SubStatusComponent
};
const login = {
  path: 'login',
  component: LoginComponent,
  outlet: 'user'
};
const update = {
  path: 'updatePass',
  component: UpdatePassComponent,
};
const routes: Routes = [
  home, node, login, update,
  { path: '**', redirectTo: '' },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { useHash: true })],
  exports: [RouterModule],
})
export class AppRouteModule { }
