create table `tatoeba_sentence` (
 `sentence_number` int not null
,`lang` varchar(3) character set ascii not null
,`text` varchar(500) not null
,`author` varchar(20) character set ascii not null
,`updated_at` datetime not null
,primary key(`sentence_number`)
);
