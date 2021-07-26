create table `app_user` (
 `id` integer primary key autoincrement
,`version` int not null
,`created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp
,`created_by` int not null
,`updated_by` int not null
,`organization_id` int not null
,`login_id` varchar(200)
,`hashed_password` varchar(200)
,`username` varchar(40)
,`role` varchar(20)
,`provider` varchar(40)
,`provider_id` varchar(40)
,`provider_access_token` text
,`provider_refresh_token` text
,`removed` tinyint(1) not null
,unique(`organization_id`, `login_id`)
);
