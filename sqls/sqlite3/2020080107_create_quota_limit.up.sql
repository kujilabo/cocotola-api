create table `quota_limit` (
 `organization_id` int not null
,`app_user_id` int not null
,`name` varchar(40) not null
,`unit` varchar(8) not null
,`date` datetime not null
,`count` int not null
,primary key(`organization_id`, `app_user_id`, `name`, `unit`, `date`)
,foreign key(`organization_id`) references `organization`(`id`)
,foreign key(`app_user_id`) references `app_user`(`id`)
);
