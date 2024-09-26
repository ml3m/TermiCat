package asciiart

// GetFrames returns the ASCII art frames of the cat animation (without whiskers)
func GetFrames() []string {
    return []string{
        `
   /\     /\  
  /  \___/  \ 
 (  o     o  )
 (     ^     )
  \_________/ 
   /       \  
  /         \ 
 /_/|     |\_\
     `,
     `
   /\     /\  
  /  \___/  \ 
 (  _     _  )
 (     ^     )
  \_________/ 
   /       \  
  /         \ 
 /_/|     |\_\
     `,
     `
 
   /\     /\  
  /  \___/  \ 
 (  _     _  )
 (     ^     )
  \_________/ 
  /         \ 
 /_/|     |\_\
      `,
    }
}

func GetDeadCat() []string {
    return []string{
     `
   /\     /\  
  /  \___/  \ 
 (  x     x  )
 (     ^     )
  \_________/ 
  /         \ 
 /_/|     |\_\
      `,
     `
   /\     /\  
  /  \___/  \ 
 (  X     X  )
 (     ^     )
  \_________/ 
  /         \ 
 /_/|     |\_\
      `,
    }
}

func GetLovelyCat() []string {
    return []string {
        `
   /\     /\  
  /  \___/  \ 
 (  o     o  )
 (     ^     )
  \_________/ 
   /       \  
  /         \ 
 /_/|     |\_\
     `,
    `
   /\     /\  
  /  \___/  \ 
 (  ♥     ♥  )
 (     ^     )
  \_________/ 
   /       \  
  /         \ 
 /_/|     |\_\
     `,
    `
              
   /\     /\  
  /  \___/  \ 
 (  ❤     ❤  )
 (     ^     )
  \_________/ 
  /         \ 
 /_/|     |\_\
     `,
    }
}

func GetLevelUpFrames() []string {
    return []string {
        `
   ★         ٭
      ⋆   ★   
 <>=========<>
 | LEVEL UP! |
 <>=========<>
     ★      ⋆ 
        ★     
  ٭           
        `,
        `
  ⋆      ★    
    ★  ⋆    ٭ 
 <>=========<>
 | LEVEL UP! |
 <>=========<>
 ٭       ★    
    ★       ⋆ 
              
        `,
        `
     ★        
 ⋆      ★     
 <>=========<>
 | LEVEL UP! |
 <>=========<>
            ٭ 
 ٭      ★     
  ★       ⋆   
        `,
    }
}

func GetFullnessFrames() []string {
    return []string {
        `
   !      !   
       !     !
 <>=========<>
 | Termi is  |
 |   Full!   |
 <>=========<>
     !   !    
  !   !    !  
        `,

    }
}
