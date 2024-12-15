package main

// Node 表示十字链表的节点
type Node struct {
	UserID      int   // 当前用户 ID
	FollowingID int   // 被关注的用户 ID
	NextRow     *Node // 行链表的下一个节点（同一个用户的关注列表）
	NextColumn  *Node // 列链表的下一个节点（同一个用户的粉丝列表）
}

// OrthogonalList 表示整个十字链表
type OrthogonalList struct {
	Rows map[int]*Node // 行索引，存储每个用户的关注链表
	Cols map[int]*Node // 列索引，存储每个用户的粉丝链表
}

// NewOrthogonalList 初始化十字链表
func NewOrthogonalList() *OrthogonalList {
	return &OrthogonalList{
		Rows: make(map[int]*Node),
		Cols: make(map[int]*Node),
	}
}

// AddFollow 添加关注关系
func (ol *OrthogonalList) AddFollow(userID, followingID int) {
	// 创建新节点
	newNode := &Node{
		UserID:      userID,
		FollowingID: followingID,
	}

	// 更新行链表（关注）
	if ol.Rows[userID] == nil {
		ol.Rows[userID] = newNode
	} else {
		current := ol.Rows[userID]
		for current.NextRow != nil {
			current = current.NextRow
		}
		current.NextRow = newNode
	}

	// 更新列链表（粉丝）
	if ol.Cols[followingID] == nil {
		ol.Cols[followingID] = newNode
	} else {
		current := ol.Cols[followingID]
		for current.NextColumn != nil {
			current = current.NextColumn
		}
		current.NextColumn = newNode
	}
}

// GetFollowings 获取用户的关注列表
func (ol *OrthogonalList) GetFollowings(userID int) []int {
	var followings []int
	current := ol.Rows[userID]
	for current != nil {
		followings = append(followings, current.FollowingID)
		current = current.NextRow
	}
	return followings
}

// GetFollowers 获取用户的粉丝列表
func (ol *OrthogonalList) GetFollowers(userID int) []int {
	var followers []int
	current := ol.Cols[userID]
	for current != nil {
		followers = append(followers, current.UserID)
		current = current.NextColumn
	}
	return followers
}

//CREATE TABLE following_x (
//id BIGINT AUTO_INCREMENT PRIMARY KEY,
//user_id BIGINT NOT NULL, -- 发起关注的用户
//following_id BIGINT NOT NULL, -- 被关注的用户
//created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//) ENGINE=InnoDB;
//CREATE TABLE follower_x (
//id BIGINT AUTO_INCREMENT PRIMARY KEY,
//user_id BIGINT NOT NULL, -- 被关注的用户
//follower_id BIGINT NOT NULL, -- 粉丝的用户
//created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
//) ENGINE=InnoDB;
