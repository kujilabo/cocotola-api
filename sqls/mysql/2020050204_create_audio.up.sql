create table `audio` (
 `id` int auto_increment
,`lang` varchar(5) character set ascii
,`text` varchar(100) character set ascii not null
,`audio_content` text character set ascii not null
,primary key(`id`)
,unique(`lang`, `text`)
);
